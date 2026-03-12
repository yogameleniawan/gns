package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/database"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionNotFound    = errors.New("session not found")
)

type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByOAuth(ctx context.Context, provider, oauthID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// User list with pagination
	ListUsers(ctx context.Context, req UserListRequest) ([]models.User, error)
	CountUsers(ctx context.Context, req UserListRequest) (int64, error)

	// Deleted users (restore)
	RestoreUser(ctx context.Context, id uuid.UUID) error
	GetDeletedUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListDeletedUsers(ctx context.Context, req UserListRequest) ([]models.User, error)
	CountDeletedUsers(ctx context.Context, req UserListRequest) (int64, error)

	// Session operations
	CreateSession(ctx context.Context, session *models.UserSession) error
	GetSessionByAccessToken(ctx context.Context, accessToken string) (*models.UserSession, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.UserSession, error)
	InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error
	InvalidateSession(ctx context.Context, sessionID uuid.UUID) error
	UpdateSession(ctx context.Context, session *models.UserSession) error
	CleanupExpiredSessions(ctx context.Context) error

	// User with roles and permissions
	GetUserWithRolesAndPermissions(ctx context.Context, userID uuid.UUID) (*models.User, []string, []string, error)
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{
		db: db,
	}
}

// --- User Operations ---

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	// Check if email already exists
	var existingUser models.User
	if err := r.db.GetDB().Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return ErrEmailAlreadyExists
	}

	if err := r.db.GetDB().WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.GetDB().WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.GetDB().WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByOAuth(ctx context.Context, provider, oauthID string) (*models.User, error) {
	var user models.User
	if err := r.db.GetDB().WithContext(ctx).
		Where("oauth_provider = ? AND oauth_id = ? AND deleted_at IS NULL", provider, oauthID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	if err := r.db.GetDB().WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if err := r.db.GetDB().WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error; err != nil {
		return err
	}
	return nil
}

// --- User List Operations ---

func (r *repository) ListUsers(ctx context.Context, req UserListRequest) ([]models.User, error) {
	var users []models.User

	query := r.db.GetDB().WithContext(ctx).Model(&models.User{}).Where("deleted_at IS NULL")

	// Apply search filter
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Apply status filter
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Apply role filter
	if req.RoleID != nil {
		query = query.Where("id IN (SELECT user_id FROM user_roles WHERE role_id = ?)", *req.RoleID)
	}

	// Apply sorting
	if req.SortBy != "" {
		allowedSortFields := map[string]bool{
			"name": true, "email": true, "created_at": true, "is_active": true,
		}
		if allowedSortFields[req.SortBy] {
			sortDir := "ASC"
			if req.SortDir == "desc" {
				sortDir = "DESC"
			}
			query = query.Order(req.SortBy + " " + sortDir)
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) CountUsers(ctx context.Context, req UserListRequest) (int64, error) {
	var count int64

	query := r.db.GetDB().WithContext(ctx).Model(&models.User{}).Where("deleted_at IS NULL")

	// Apply search filter
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Apply status filter
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Apply role filter
	if req.RoleID != nil {
		query = query.Where("id IN (SELECT user_id FROM user_roles WHERE role_id = ?)", *req.RoleID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// --- Deleted Users (Restore) ---

func (r *repository) RestoreUser(ctx context.Context, id uuid.UUID) error {
	if err := r.db.GetDB().WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetDeletedUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.GetDB().WithContext(ctx).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) ListDeletedUsers(ctx context.Context, req UserListRequest) ([]models.User, error) {
	var users []models.User

	query := r.db.GetDB().WithContext(ctx).Model(&models.User{}).Where("deleted_at IS NOT NULL")

	// Apply search filter
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Apply sorting
	if req.SortBy != "" {
		allowedSortFields := map[string]bool{
			"name": true, "email": true, "created_at": true, "deleted_at": true,
		}
		if allowedSortFields[req.SortBy] {
			sortDir := "ASC"
			if req.SortDir == "desc" {
				sortDir = "DESC"
			}
			query = query.Order(req.SortBy + " " + sortDir)
		}
	} else {
		query = query.Order("deleted_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) CountDeletedUsers(ctx context.Context, req UserListRequest) (int64, error) {
	var count int64

	query := r.db.GetDB().WithContext(ctx).Model(&models.User{}).Where("deleted_at IS NOT NULL")

	// Apply search filter
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// --- Session Operations ---

func (r *repository) CreateSession(ctx context.Context, session *models.UserSession) error {
	if err := r.db.GetDB().WithContext(ctx).Create(session).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetSessionByAccessToken(ctx context.Context, accessToken string) (*models.UserSession, error) {
	var session models.UserSession
	if err := r.db.GetDB().WithContext(ctx).
		Where("access_token = ? AND is_active = ?", accessToken, true).
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *repository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.UserSession, error) {
	var session models.UserSession
	if err := r.db.GetDB().WithContext(ctx).
		Where("refresh_token = ? AND is_active = ?", refreshToken, true).
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *repository) InvalidateUserSessions(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.GetDB().WithContext(ctx).
		Model(&models.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) InvalidateSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := r.db.GetDB().WithContext(ctx).
		Model(&models.UserSession{}).
		Where("id = ?", sessionID).
		Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateSession(ctx context.Context, session *models.UserSession) error {
	session.UpdatedAt = time.Now()
	if err := r.db.GetDB().WithContext(ctx).Save(session).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) CleanupExpiredSessions(ctx context.Context) error {
	now := time.Now()
	if err := r.db.GetDB().WithContext(ctx).
		Model(&models.UserSession{}).
		Where("refresh_token_expires_at < ? OR is_active = ?", now, false).
		Delete(&models.UserSession{}).Error; err != nil {
		return err
	}
	return nil
}

// --- User with Roles and Permissions ---

func (r *repository) GetUserWithRolesAndPermissions(ctx context.Context, userID uuid.UUID) (*models.User, []string, []string, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get roles
	var roles []string
	if err := r.db.GetDB().WithContext(ctx).
		Table("user_roles").
		Select("roles.name").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Scan(&roles).Error; err != nil {
		return nil, nil, nil, err
	}

	// Get permissions
	var permissions []string
	if err := r.db.GetDB().WithContext(ctx).
		Table("user_roles").
		Select("DISTINCT permissions.name").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ?", userID).
		Scan(&permissions).Error; err != nil {
		return nil, nil, nil, err
	}

	return user, roles, permissions, nil
}
