package models

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role in the RBAC system
type Role struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description *string   `gorm:"type:text"`
	IsSystem    bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission represents a permission in the RBAC system
type Permission struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Module      string    `gorm:"type:varchar(100);not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	RoleID       int       `gorm:"not null"`
	PermissionID int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	RoleID    int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// ModuleAccess represents module-level access control for roles
type ModuleAccess struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	RoleID     int       `gorm:"not null"`
	ModuleName string    `gorm:"type:varchar(100);not null"`
	CanView    bool      `gorm:"default:false"`
	CanCreate  bool      `gorm:"default:false"`
	CanEdit    bool      `gorm:"default:false"`
	CanDelete  bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"not null;default:now()"`
	UpdatedAt  time.Time `gorm:"not null;default:now()"`
}

func (ModuleAccess) TableName() string {
	return "module_access"
}
