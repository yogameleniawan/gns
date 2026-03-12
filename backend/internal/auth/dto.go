package auth

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type GoogleOAuthRequest struct {
	IdToken string `json:"id_token" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name" validate:"omitempty,min=2"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// Response DTOs

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	AvatarURL     string    `json:"avatar_url,omitempty"`
	IsOAuth       bool      `json:"is_oauth"`
	OAuthProvider string    `json:"oauth_provider,omitempty"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	Roles         []string  `json:"roles,omitempty"`
	Permissions   []string  `json:"permissions,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProfileResponse struct {
	User UserResponse `json:"user"`
}

// Google OAuth Response from Google API
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

// User Management DTOs

// UserListRequest for paginated user list with filters
type UserListRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
	IsActive *bool  `json:"is_active"`
	RoleID   *int   `json:"role_id"`
}

// UserListResponse for paginated user list
type UserListResponse struct {
	Items      []UserResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// CreateUserRequest for admin creating new user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
	RoleIDs  []int  `json:"role_ids"`
}

// AdminUpdateUserRequest for admin updating user
type AdminUpdateUserRequest struct {
	Name      string `json:"name" validate:"omitempty,min=2"`
	Email     string `json:"email" validate:"omitempty,email"`
	IsActive  *bool  `json:"is_active"`
	AvatarURL string `json:"avatar_url"`
	RoleIDs   []int  `json:"role_ids"`
}
