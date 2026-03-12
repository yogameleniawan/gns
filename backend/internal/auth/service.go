package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/cache"
	"github.com/base-go/backend/pkg/config"
	"github.com/base-go/backend/pkg/token"
	"github.com/base-go/backend/pkg/utils"
	"github.com/base-go/backend/pkg/validator"
)

var (
	ErrWeakPassword     = errors.New("password must be at least 8 characters")
	ErrUserInactive     = errors.New("user account is inactive")
	ErrPasswordMismatch = errors.New("old password is incorrect")
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, int, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, int, error)
	LoginWithGoogle(ctx context.Context, req GoogleOAuthRequest) (*AuthResponse, int, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, int, error)
	Logout(ctx context.Context, userID uuid.UUID) (int, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, int, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*UserResponse, int, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req ChangePasswordRequest) (int, error)

	// User management (admin)
	ListUsers(ctx context.Context, req UserListRequest) (*UserListResponse, int, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, int, error)
	CreateUserByAdmin(ctx context.Context, req CreateUserRequest) (*UserResponse, int, error)
	UpdateUserByAdmin(ctx context.Context, id uuid.UUID, req AdminUpdateUserRequest) (*UserResponse, int, error)
	DeleteUserByAdmin(ctx context.Context, id uuid.UUID) (int, error)
	ToggleUserStatus(ctx context.Context, id uuid.UUID) (*UserResponse, int, error)

	// Deleted users (restore)
	ListDeletedUsers(ctx context.Context, req UserListRequest) (*UserListResponse, int, error)
	RestoreUserByAdmin(ctx context.Context, id uuid.UUID) (*UserResponse, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, int, error) {
	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Validate password strength
	if len(req.Password) < 8 {
		return nil, http.StatusBadRequest, ErrWeakPassword
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return nil, http.StatusInternalServerError, errors.New("failed to process password")
	}

	// Create user
	user := &models.User{
		Email:         req.Email,
		PasswordHash:  &hashedPassword,
		Name:          req.Name,
		IsOAuth:       false,
		IsActive:      true,
		EmailVerified: false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return nil, http.StatusConflict, err
		}
		logrus.WithError(err).Error("Failed to create user")
		return nil, http.StatusInternalServerError, errors.New("failed to create user")
	}

	// Assign default "User" role
	if err := s.assignDefaultRole(ctx, user.ID); err != nil {
		logrus.WithError(err).Error("Failed to assign default role")
		// Continue anyway, user is created
	}

	// Generate tokens
	authResp, err := s.generateAuthResponse(ctx, user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return authResp, http.StatusCreated, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, int, error) {
	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, http.StatusUnauthorized, ErrInvalidCredentials
		}
		logrus.WithError(err).Error("Failed to get user")
		return nil, http.StatusInternalServerError, errors.New("failed to authenticate")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, http.StatusForbidden, ErrUserInactive
	}

	// Check if OAuth user
	if user.IsOAuth {
		return nil, http.StatusBadRequest, errors.New("please use OAuth login")
	}

	// Verify password
	if user.PasswordHash == nil || !utils.CheckPassword(req.Password, *user.PasswordHash) {
		return nil, http.StatusUnauthorized, ErrInvalidCredentials
	}

	// Invalidate old sessions
	if err := s.repo.InvalidateUserSessions(ctx, user.ID); err != nil {
		logrus.WithError(err).Warn("Failed to invalidate old sessions")
	}

	// Generate tokens
	authResp, err := s.generateAuthResponse(ctx, user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return authResp, http.StatusOK, nil
}

func (s *service) LoginWithGoogle(ctx context.Context, req GoogleOAuthRequest) (*AuthResponse, int, error) {
	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Verify Google ID token
	googleUser, err := utils.VerifyGoogleIDToken(ctx, req.IdToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to verify Google token")
		return nil, http.StatusUnauthorized, errors.New("invalid Google token")
	}

	// Try to find existing user by OAuth
	user, err := s.repo.GetUserByOAuth(ctx, "google", googleUser.Sub)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		logrus.WithError(err).Error("Failed to get user by OAuth")
		return nil, http.StatusInternalServerError, errors.New("failed to authenticate")
	}

	// If user not found by OAuth, try by email
	if user == nil {
		user, err = s.repo.GetUserByEmail(ctx, googleUser.Email)
		if err != nil && !errors.Is(err, ErrUserNotFound) {
			logrus.WithError(err).Error("Failed to get user by email")
			return nil, http.StatusInternalServerError, errors.New("failed to authenticate")
		}
	}

	// Create new user if not exists
	if user == nil {
		emailVerified := googleUser.EmailVerified == "true"
		user = &models.User{
			Email:         googleUser.Email,
			Name:          googleUser.Name,
			AvatarURL:     &googleUser.Picture,
			IsOAuth:       true,
			OAuthProvider: stringPtr("google"),
			OAuthID:       &googleUser.Sub,
			IsActive:      true,
			EmailVerified: emailVerified,
		}

		if err := s.repo.CreateUser(ctx, user); err != nil {
			if errors.Is(err, ErrEmailAlreadyExists) {
				return nil, http.StatusConflict, err
			}
			logrus.WithError(err).Error("Failed to create OAuth user")
			return nil, http.StatusInternalServerError, errors.New("failed to create user")
		}

		// Assign default role
		if err := s.assignDefaultRole(ctx, user.ID); err != nil {
			logrus.WithError(err).Error("Failed to assign default role")
		}
	} else {
		// Update OAuth info if user exists but wasn't OAuth before
		if !user.IsOAuth {
			user.IsOAuth = true
			user.OAuthProvider = stringPtr("google")
			user.OAuthID = &googleUser.Sub
			user.AvatarURL = &googleUser.Picture
			if err := s.repo.UpdateUser(ctx, user); err != nil {
				logrus.WithError(err).Warn("Failed to update user OAuth info")
			}
		}
	}

	// Check if user is active
	if !user.IsActive {
		return nil, http.StatusForbidden, ErrUserInactive
	}

	// Invalidate old sessions
	if err := s.repo.InvalidateUserSessions(ctx, user.ID); err != nil {
		logrus.WithError(err).Warn("Failed to invalidate old sessions")
	}

	// Generate tokens
	authResp, err := s.generateAuthResponse(ctx, user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return authResp, http.StatusOK, nil
}

func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, int, error) {
	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Get session by refresh token
	session, err := s.repo.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, http.StatusUnauthorized, errors.New("invalid refresh token")
		}
		logrus.WithError(err).Error("Failed to get session")
		return nil, http.StatusInternalServerError, errors.New("failed to refresh token")
	}

	// Check if refresh token is expired
	if time.Now().After(session.RefreshTokenExpiresAt) {
		return nil, http.StatusUnauthorized, errors.New("refresh token expired")
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user")
		return nil, http.StatusInternalServerError, errors.New("failed to refresh token")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, http.StatusForbidden, ErrUserInactive
	}

	// Invalidate old session
	if err := s.repo.InvalidateSession(ctx, session.ID); err != nil {
		logrus.WithError(err).Warn("Failed to invalidate old session")
	}

	// Generate new tokens
	authResp, err := s.generateAuthResponse(ctx, user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return authResp, http.StatusOK, nil
}

func (s *service) Logout(ctx context.Context, userID uuid.UUID) (int, error) {
	// Get user to get email for cache key
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user for logout")
		// Continue with logout even if user not found
	}

	// Invalidate all user sessions
	if err := s.repo.InvalidateUserSessions(ctx, userID); err != nil {
		logrus.WithError(err).Error("Failed to invalidate sessions")
		return http.StatusInternalServerError, errors.New("failed to logout")
	}

	// Delete tokens from Redis cache
	if user != nil {
		ch := cache.NewCache()
		accessTokenKey := fmt.Sprintf("%s_USER_ACCESS_TOKEN", user.Email)
		refreshTokenKey := fmt.Sprintf("%s_USER_REFRESH_TOKEN", user.Email)

		if err := ch.Delete(ctx, accessTokenKey); err != nil {
			logrus.WithError(err).Error("Failed to delete access token from cache")
		}

		if err := ch.Delete(ctx, refreshTokenKey); err != nil {
			logrus.WithError(err).Error("Failed to delete refresh token from cache")
		}
	}

	return http.StatusOK, nil
}

func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, int, error) {
	// Get user with roles and permissions
	user, roles, permissions, err := s.repo.GetUserWithRolesAndPermissions(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get user profile")
		return nil, http.StatusInternalServerError, errors.New("failed to get profile")
	}

	userResp := s.mapUserToResponse(user, roles, permissions)

	return &ProfileResponse{
		User: userResp,
	}, http.StatusOK, nil
}

func (s *service) UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*UserResponse, int, error) {
	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get user")
		return nil, http.StatusInternalServerError, errors.New("failed to update profile")
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	// Save user
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		logrus.WithError(err).Error("Failed to update user")
		return nil, http.StatusInternalServerError, errors.New("failed to update profile")
	}

	// Get updated user with roles and permissions
	user, roles, permissions, err := s.repo.GetUserWithRolesAndPermissions(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get updated user")
		return nil, http.StatusInternalServerError, errors.New("failed to get updated profile")
	}

	userResp := s.mapUserToResponse(user, roles, permissions)

	return &userResp, http.StatusOK, nil
}

func (s *service) ChangePassword(ctx context.Context, userID uuid.UUID, req ChangePasswordRequest) (int, error) {
	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		return http.StatusBadRequest, err
	}

	// Validate new password strength
	if len(req.NewPassword) < 8 {
		return http.StatusBadRequest, ErrWeakPassword
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get user")
		return http.StatusInternalServerError, errors.New("failed to change password")
	}

	// Check if OAuth user
	if user.IsOAuth {
		return http.StatusBadRequest, errors.New("cannot change password for OAuth users")
	}

	// Verify old password
	if user.PasswordHash == nil || !utils.CheckPassword(req.OldPassword, *user.PasswordHash) {
		return http.StatusUnauthorized, ErrPasswordMismatch
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return http.StatusInternalServerError, errors.New("failed to process password")
	}

	// Update password
	user.PasswordHash = &hashedPassword
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		logrus.WithError(err).Error("Failed to update password")
		return http.StatusInternalServerError, errors.New("failed to change password")
	}

	// Invalidate all sessions (force re-login)
	if err := s.repo.InvalidateUserSessions(ctx, userID); err != nil {
		logrus.WithError(err).Warn("Failed to invalidate sessions")
	}

	return http.StatusOK, nil
}

// --- Helper Functions ---

func (s *service) generateAuthResponse(ctx context.Context, user *models.User) (*AuthResponse, error) {
	cfg := config.GetConfig()

	// Get user roles and permissions
	_, roles, permissions, err := s.repo.GetUserWithRolesAndPermissions(ctx, user.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user roles and permissions")
		roles = []string{}
		permissions = []string{}
	}

	// Generate JWT tokens
	rolesStr := ""
	if len(roles) > 0 {
		rolesStr = roles[0] // Use first role for token
	}

	accessToken, refreshToken, err := token.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		rolesStr,
		"", // mitraType not used in this system
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate tokens")
		return nil, errors.New("failed to generate tokens")
	}

	// Calculate expiration times
	now := time.Now()
	accessExp := now.Add(time.Duration(cfg.Auth.TokenExpiration) * time.Second)
	refreshExp := now.Add(time.Duration(cfg.Auth.RefreshTokenExpiration) * time.Second)

	// Create session
	session := &models.UserSession{
		UserID:                user.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExp,
		RefreshTokenExpiresAt: refreshExp,
		IsActive:              true,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		logrus.WithError(err).Error("Failed to create session")
		return nil, errors.New("failed to create session")
	}

	// Store tokens in Redis cache for validation
	ch := cache.NewCache()
	accessTokenKey := fmt.Sprintf("%s_USER_ACCESS_TOKEN", user.Email)
	refreshTokenKey := fmt.Sprintf("%s_USER_REFRESH_TOKEN", user.Email)

	// Store access token with expiration
	if err := ch.Set(ctx, accessTokenKey, accessToken, time.Duration(cfg.Auth.TokenExpiration)*time.Second); err != nil {
		logrus.WithError(err).Error("Failed to cache access token")
		// Don't fail the request, just log the error
	}

	// Store refresh token with expiration
	if err := ch.Set(ctx, refreshTokenKey, refreshToken, time.Duration(cfg.Auth.RefreshTokenExpiration)*time.Second); err != nil {
		logrus.WithError(err).Error("Failed to cache refresh token")
		// Don't fail the request, just log the error
	}

	// Build response
	userResp := s.mapUserToResponse(user, roles, permissions)

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    cfg.Auth.TokenExpiration,
		User:         userResp,
	}, nil
}

func (s *service) mapUserToResponse(user *models.User, roles, permissions []string) UserResponse {
	avatarURL := ""
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	oauthProvider := ""
	if user.OAuthProvider != nil {
		oauthProvider = *user.OAuthProvider
	}

	return UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		AvatarURL:     avatarURL,
		IsOAuth:       user.IsOAuth,
		OAuthProvider: oauthProvider,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
		Roles:         roles,
		Permissions:   permissions,
		CreatedAt:     user.CreatedAt,
	}
}

func (s *service) assignDefaultRole(ctx context.Context, userID uuid.UUID) error {
	// This will be implemented when we create the RBAC module
	// For now, we'll insert directly
	db := s.repo.(*repository).db.GetDB()

	// Get "User" role ID
	var roleID int
	if err := db.Table("roles").
		Select("id").
		Where("name = ?", "User").
		Scan(&roleID).Error; err != nil {
		return err
	}

	// Assign role to user
	if err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		userID, roleID).Error; err != nil {
		return err
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}

// --- User Management (Admin) ---

func (s *service) ListUsers(ctx context.Context, req UserListRequest) (*UserListResponse, int, error) {
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// Get users
	users, err := s.repo.ListUsers(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to list users")
		return nil, http.StatusInternalServerError, errors.New("failed to list users")
	}

	// Get total count
	total, err := s.repo.CountUsers(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to count users")
		return nil, http.StatusInternalServerError, errors.New("failed to count users")
	}

	// Map users to response with roles
	items := make([]UserResponse, len(users))
	for i, user := range users {
		_, roles, permissions, _ := s.repo.GetUserWithRolesAndPermissions(ctx, user.ID)
		items[i] = s.mapUserToResponse(&user, roles, permissions)
	}

	// Calculate total pages
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &UserListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, http.StatusOK, nil
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, int, error) {
	user, roles, permissions, err := s.repo.GetUserWithRolesAndPermissions(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get user")
		return nil, http.StatusInternalServerError, errors.New("failed to get user")
	}

	response := s.mapUserToResponse(user, roles, permissions)
	return &response, http.StatusOK, nil
}

func (s *service) CreateUserByAdmin(ctx context.Context, req CreateUserRequest) (*UserResponse, int, error) {
	// Validate request
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return nil, http.StatusInternalServerError, errors.New("failed to process password")
	}

	// Create user
	user := &models.User{
		Email:         req.Email,
		PasswordHash:  &hashedPassword,
		Name:          req.Name,
		IsOAuth:       false,
		IsActive:      true,
		EmailVerified: false,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return nil, http.StatusConflict, err
		}
		logrus.WithError(err).Error("Failed to create user")
		return nil, http.StatusInternalServerError, errors.New("failed to create user")
	}

	// Assign roles if provided
	if len(req.RoleIDs) > 0 {
		db := s.repo.(*repository).db.GetDB()
		for _, roleID := range req.RoleIDs {
			if err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				user.ID, roleID).Error; err != nil {
				logrus.WithError(err).Error("Failed to assign role")
			}
		}
	} else {
		// Assign default role
		if err := s.assignDefaultRole(ctx, user.ID); err != nil {
			logrus.WithError(err).Error("Failed to assign default role")
		}
	}

	// Get user with roles
	_, roles, permissions, _ := s.repo.GetUserWithRolesAndPermissions(ctx, user.ID)
	response := s.mapUserToResponse(user, roles, permissions)

	return &response, http.StatusCreated, nil
}

func (s *service) UpdateUserByAdmin(ctx context.Context, id uuid.UUID, req AdminUpdateUserRequest) (*UserResponse, int, error) {
	// Get existing user
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get user")
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		// Check if new email exists
		existing, _ := s.repo.GetUserByEmail(ctx, req.Email)
		if existing != nil && existing.ID != user.ID {
			return nil, http.StatusConflict, ErrEmailAlreadyExists
		}
		user.Email = req.Email
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	// Save user
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		logrus.WithError(err).Error("Failed to update user")
		return nil, http.StatusInternalServerError, errors.New("failed to update user")
	}

	// Update roles if provided
	if req.RoleIDs != nil {
		db := s.repo.(*repository).db.GetDB()
		// Remove existing roles
		if err := db.Exec("DELETE FROM user_roles WHERE user_id = ?", user.ID).Error; err != nil {
			logrus.WithError(err).Error("Failed to remove existing roles")
		}
		// Add new roles
		for _, roleID := range req.RoleIDs {
			if err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				user.ID, roleID).Error; err != nil {
				logrus.WithError(err).Error("Failed to assign role")
			}
		}
	}

	// Get user with updated roles
	_, roles, permissions, _ := s.repo.GetUserWithRolesAndPermissions(ctx, user.ID)
	response := s.mapUserToResponse(user, roles, permissions)

	return &response, http.StatusOK, nil
}

func (s *service) DeleteUserByAdmin(ctx context.Context, id uuid.UUID) (int, error) {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, errors.New("failed to get user")
	}

	// Soft delete
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		logrus.WithError(err).Error("Failed to delete user")
		return http.StatusInternalServerError, errors.New("failed to delete user")
	}

	// Invalidate all user sessions
	if err := s.repo.InvalidateUserSessions(ctx, id); err != nil {
		logrus.WithError(err).Error("Failed to invalidate user sessions")
	}

	return http.StatusOK, nil
}

func (s *service) ToggleUserStatus(ctx context.Context, id uuid.UUID) (*UserResponse, int, error) {
	// Get user
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get user")
	}

	// Toggle status
	user.IsActive = !user.IsActive

	// Save
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		logrus.WithError(err).Error("Failed to update user status")
		return nil, http.StatusInternalServerError, errors.New("failed to update user status")
	}

	// If user is now inactive, invalidate sessions
	if !user.IsActive {
		if err := s.repo.InvalidateUserSessions(ctx, id); err != nil {
			logrus.WithError(err).Error("Failed to invalidate user sessions")
		}
	}

	_, roles, permissions, _ := s.repo.GetUserWithRolesAndPermissions(ctx, user.ID)
	response := s.mapUserToResponse(user, roles, permissions)

	return &response, http.StatusOK, nil
}

// --- Deleted Users (Restore) ---

func (s *service) ListDeletedUsers(ctx context.Context, req UserListRequest) (*UserListResponse, int, error) {
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// Get deleted users
	users, err := s.repo.ListDeletedUsers(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to list deleted users")
		return nil, http.StatusInternalServerError, errors.New("failed to list deleted users")
	}

	// Get total count
	total, err := s.repo.CountDeletedUsers(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to count deleted users")
		return nil, http.StatusInternalServerError, errors.New("failed to count deleted users")
	}

	// Map users to response
	items := make([]UserResponse, len(users))
	for i, user := range users {
		items[i] = s.mapUserToResponse(&user, nil, nil)
	}

	// Calculate total pages
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &UserListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, http.StatusOK, nil
}

func (s *service) RestoreUserByAdmin(ctx context.Context, id uuid.UUID) (*UserResponse, int, error) {
	// Check if deleted user exists
	user, err := s.repo.GetDeletedUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, http.StatusNotFound, errors.New("deleted user not found")
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get deleted user")
	}

	// Restore user
	if err := s.repo.RestoreUser(ctx, id); err != nil {
		logrus.WithError(err).Error("Failed to restore user")
		return nil, http.StatusInternalServerError, errors.New("failed to restore user")
	}

	// Get fresh user data
	user, err = s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("failed to get restored user")
	}

	_, roles, permissions, _ := s.repo.GetUserWithRolesAndPermissions(ctx, user.ID)
	response := s.mapUserToResponse(user, roles, permissions)

	return &response, http.StatusOK, nil
}
