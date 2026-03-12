package rbac

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Description string `json:"description"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Module      string `json:"module" validate:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Module      string `json:"module" validate:"omitempty,min=2,max=100"`
	Description string `json:"description"`
}

type AssignPermissionsToRoleRequest struct {
	PermissionIDs []int `json:"permission_ids" validate:"required,min=1"`
}

type AssignRolesToUserRequest struct {
	RoleIDs []int `json:"role_ids" validate:"required,min=1"`
}

type UpdateModuleAccessRequest struct {
	ModuleName string `json:"module_name" validate:"required"`
	CanView    bool   `json:"can_view"`
	CanCreate  bool   `json:"can_create"`
	CanEdit    bool   `json:"can_edit"`
	CanDelete  bool   `json:"can_delete"`
}

type CheckPermissionRequest struct {
	Permission string `json:"permission" validate:"required"`
}

type CheckModuleAccessRequest struct {
	Module string `json:"module" validate:"required"`
	Action string `json:"action" validate:"required,oneof=view create edit delete"`
}

// Response DTOs

type RoleResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleWithPermissionsResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	IsSystem    bool                 `json:"is_system"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type PermissionResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Module      string    `json:"module"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PermissionsByModuleResponse struct {
	Module      string               `json:"module"`
	Permissions []PermissionResponse `json:"permissions"`
}

type ModuleAccessResponse struct {
	ID         int       `json:"id"`
	RoleID     int       `json:"role_id"`
	ModuleName string    `json:"module_name"`
	CanView    bool      `json:"can_view"`
	CanCreate  bool      `json:"can_create"`
	CanEdit    bool      `json:"can_edit"`
	CanDelete  bool      `json:"can_delete"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRoleResponse struct {
	UserID    uuid.UUID      `json:"user_id"`
	UserEmail string         `json:"user_email"`
	UserName  string         `json:"user_name"`
	Roles     []RoleResponse `json:"roles"`
}

type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

type CheckModuleAccessResponse struct {
	HasAccess bool `json:"has_access"`
}
