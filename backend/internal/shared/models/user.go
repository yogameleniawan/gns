package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash  *string    `gorm:"type:varchar(255)"`
	Name          string     `gorm:"type:varchar(255);not null"`
	AvatarURL     *string    `gorm:"type:text"`
	IsOAuth       bool       `gorm:"column:is_oauth;default:false"`
	OAuthProvider *string    `gorm:"column:oauth_provider;type:varchar(50)"`
	OAuthID       *string    `gorm:"column:oauth_id;type:varchar(255)"`
	IsActive      bool       `gorm:"default:true"`
	EmailVerified bool       `gorm:"default:false"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()"`
	DeletedAt     *time.Time `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

// UserSession represents a user session for managing active sessions
type UserSession struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID                uuid.UUID `gorm:"type:uuid;not null"`
	AccessToken           string    `gorm:"type:text;not null"`
	RefreshToken          string    `gorm:"type:text;not null"`
	AccessTokenExpiresAt  time.Time `gorm:"not null"`
	RefreshTokenExpiresAt time.Time `gorm:"not null"`
	UserAgent             *string   `gorm:"type:text"`
	IPAddress             *string   `gorm:"type:varchar(45)"`
	IsActive              bool      `gorm:"default:true"`
	CreatedAt             time.Time `gorm:"not null;default:now()"`
	UpdatedAt             time.Time `gorm:"not null;default:now()"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
