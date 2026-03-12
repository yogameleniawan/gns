package rbac

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/validator"
)

type Service interface {
	// Role operations
	CreateRole(ctx context.Context, req CreateRoleRequest) (*RoleResponse, int, error)
	GetAllRoles(ctx context.Context) ([]RoleResponse, int, error)
	GetRoleByID(ctx context.Context, id int) (*RoleWithPermissionsResponse, int, error)
	UpdateRole(ctx context.Context, id int, req UpdateRoleRequest) (*RoleResponse, int, error)
	DeleteRole(ctx context.Context, id int) (int, error)

	// Permission operations
	CreatePermission(ctx context.Context, req CreatePermissionRequest) (*PermissionResponse, int, error)
	GetAllPermissions(ctx context.Context) ([]PermissionResponse, int, error)
	GetPermissionsByModule(ctx context.Context) ([]PermissionsByModuleResponse, int, error)
	GetPermissionByID(ctx context.Context, id int) (*PermissionResponse, int, error)
	UpdatePermission(ctx context.Context, id int, req UpdatePermissionRequest) (*PermissionResponse, int, error)
	DeletePermission(ctx context.Context, id int) (int, error)

	// Role-Permission operations
	AssignPermissionsToRole(ctx context.Context, roleID int, req AssignPermissionsToRoleRequest) (int, error)
	GetRolePermissions(ctx context.Context, roleID int) ([]PermissionResponse, int, error)

	// User-Role operations
	AssignRolesToUser(ctx context.Context, userID uuid.UUID, req AssignRolesToUserRequest) (int, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]RoleResponse, int, error)

	// Module access operations
	UpdateModuleAccess(ctx context.Context, roleID int, req UpdateModuleAccessRequest) (*ModuleAccessResponse, int, error)
	GetModuleAccessByRole(ctx context.Context, roleID int) ([]ModuleAccessResponse, int, error)

	// Permission checking
	CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, int, error)
	CheckModuleAccess(ctx context.Context, userID uuid.UUID, module string, action string) (bool, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// --- Role Operations ---

func (s *service) CreateRole(ctx context.Context, req CreateRoleRequest) (*RoleResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	role := &models.Role{
		Name:        req.Name,
		Description: stringPtr(req.Description),
		IsSystem:    false,
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		if errors.Is(err, ErrRoleAlreadyExists) {
			return nil, http.StatusConflict, err
		}
		logrus.WithError(err).Error("Failed to create role")
		return nil, http.StatusInternalServerError, errors.New("failed to create role")
	}

	resp := s.mapRoleToResponse(role)
	return &resp, http.StatusCreated, nil
}

func (s *service) GetAllRoles(ctx context.Context) ([]RoleResponse, int, error) {
	roles, err := s.repo.GetAllRoles(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to get roles")
		return nil, http.StatusInternalServerError, errors.New("failed to get roles")
	}

	responses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = s.mapRoleToResponse(&role)
	}

	return responses, http.StatusOK, nil
}

func (s *service) GetRoleByID(ctx context.Context, id int) (*RoleWithPermissionsResponse, int, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get role")
		return nil, http.StatusInternalServerError, errors.New("failed to get role")
	}

	permissions, err := s.repo.GetRolePermissions(ctx, id)
	if err != nil {
		logrus.WithError(err).Error("Failed to get role permissions")
		permissions = []models.Permission{}
	}

	return s.mapRoleWithPermissionsToResponse(role, permissions), http.StatusOK, nil
}

func (s *service) UpdateRole(ctx context.Context, id int, req UpdateRoleRequest) (*RoleResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get role")
		return nil, http.StatusInternalServerError, errors.New("failed to get role")
	}

	if role.IsSystem {
		return nil, http.StatusForbidden, ErrSystemRoleProtected
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = stringPtr(req.Description)
	}
	role.UpdatedAt = time.Now()

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		logrus.WithError(err).Error("Failed to update role")
		return nil, http.StatusInternalServerError, errors.New("failed to update role")
	}

	resp := s.mapRoleToResponse(role)
	return &resp, http.StatusOK, nil
}

func (s *service) DeleteRole(ctx context.Context, id int) (int, error) {
	if err := s.repo.DeleteRole(ctx, id); err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return http.StatusNotFound, err
		}
		if errors.Is(err, ErrSystemRoleProtected) {
			return http.StatusForbidden, err
		}
		logrus.WithError(err).Error("Failed to delete role")
		return http.StatusInternalServerError, errors.New("failed to delete role")
	}

	return http.StatusOK, nil
}

// --- Permission Operations ---

func (s *service) CreatePermission(ctx context.Context, req CreatePermissionRequest) (*PermissionResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	permission := &models.Permission{
		Name:        req.Name,
		Module:      req.Module,
		Description: stringPtr(req.Description),
	}

	if err := s.repo.CreatePermission(ctx, permission); err != nil {
		logrus.WithError(err).Error("Failed to create permission")
		return nil, http.StatusInternalServerError, errors.New("failed to create permission")
	}

	resp := s.mapPermissionToResponse(permission)
	return &resp, http.StatusCreated, nil
}

func (s *service) GetAllPermissions(ctx context.Context) ([]PermissionResponse, int, error) {
	permissions, err := s.repo.GetAllPermissions(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to get permissions")
		return nil, http.StatusInternalServerError, errors.New("failed to get permissions")
	}

	responses := make([]PermissionResponse, len(permissions))
	for i, perm := range permissions {
		responses[i] = s.mapPermissionToResponse(&perm)
	}

	return responses, http.StatusOK, nil
}

func (s *service) GetPermissionsByModule(ctx context.Context) ([]PermissionsByModuleResponse, int, error) {
	grouped, err := s.repo.GetPermissionsByModuleGrouped(ctx)
	if err != nil {
		logrus.WithError(err).Error("Failed to get permissions by module")
		return nil, http.StatusInternalServerError, errors.New("failed to get permissions by module")
	}

	responses := make([]PermissionsByModuleResponse, 0, len(grouped))
	for module, perms := range grouped {
		permResponses := make([]PermissionResponse, len(perms))
		for i, perm := range perms {
			permResponses[i] = s.mapPermissionToResponse(&perm)
		}
		responses = append(responses, PermissionsByModuleResponse{
			Module:      module,
			Permissions: permResponses,
		})
	}

	return responses, http.StatusOK, nil
}

func (s *service) GetPermissionByID(ctx context.Context, id int) (*PermissionResponse, int, error) {
	permission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get permission")
		return nil, http.StatusInternalServerError, errors.New("failed to get permission")
	}

	resp := s.mapPermissionToResponse(permission)
	return &resp, http.StatusOK, nil
}

func (s *service) UpdatePermission(ctx context.Context, id int, req UpdatePermissionRequest) (*PermissionResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	permission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			return nil, http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to get permission")
		return nil, http.StatusInternalServerError, errors.New("failed to get permission")
	}

	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Module != "" {
		permission.Module = req.Module
	}
	if req.Description != "" {
		permission.Description = stringPtr(req.Description)
	}
	permission.UpdatedAt = time.Now()

	if err := s.repo.UpdatePermission(ctx, permission); err != nil {
		logrus.WithError(err).Error("Failed to update permission")
		return nil, http.StatusInternalServerError, errors.New("failed to update permission")
	}

	resp := s.mapPermissionToResponse(permission)
	return &resp, http.StatusOK, nil
}

func (s *service) DeletePermission(ctx context.Context, id int) (int, error) {
	if err := s.repo.DeletePermission(ctx, id); err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			return http.StatusNotFound, err
		}
		logrus.WithError(err).Error("Failed to delete permission")
		return http.StatusInternalServerError, errors.New("failed to delete permission")
	}

	return http.StatusOK, nil
}

// --- Role-Permission Operations ---

func (s *service) AssignPermissionsToRole(ctx context.Context, roleID int, req AssignPermissionsToRoleRequest) (int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return http.StatusBadRequest, err
	}

	// Check if role exists
	if _, err := s.repo.GetRoleByID(ctx, roleID); err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, errors.New("failed to get role")
	}

	if err := s.repo.AssignPermissionsToRole(ctx, roleID, req.PermissionIDs); err != nil {
		logrus.WithError(err).Error("Failed to assign permissions to role")
		return http.StatusInternalServerError, errors.New("failed to assign permissions")
	}

	return http.StatusOK, nil
}

func (s *service) GetRolePermissions(ctx context.Context, roleID int) ([]PermissionResponse, int, error) {
	permissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get role permissions")
		return nil, http.StatusInternalServerError, errors.New("failed to get role permissions")
	}

	responses := make([]PermissionResponse, len(permissions))
	for i, perm := range permissions {
		responses[i] = s.mapPermissionToResponse(&perm)
	}

	return responses, http.StatusOK, nil
}

// --- User-Role Operations ---

func (s *service) AssignRolesToUser(ctx context.Context, userID uuid.UUID, req AssignRolesToUserRequest) (int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return http.StatusBadRequest, err
	}

	if err := s.repo.AssignRolesToUser(ctx, userID, req.RoleIDs); err != nil {
		logrus.WithError(err).Error("Failed to assign roles to user")
		return http.StatusInternalServerError, errors.New("failed to assign roles")
	}

	return http.StatusOK, nil
}

func (s *service) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]RoleResponse, int, error) {
	roles, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user roles")
		return nil, http.StatusInternalServerError, errors.New("failed to get user roles")
	}

	responses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = s.mapRoleToResponse(&role)
	}

	return responses, http.StatusOK, nil
}

// --- Module Access Operations ---

func (s *service) UpdateModuleAccess(ctx context.Context, roleID int, req UpdateModuleAccessRequest) (*ModuleAccessResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Check if role exists
	if _, err := s.repo.GetRoleByID(ctx, roleID); err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, errors.New("failed to get role")
	}

	moduleAccess := &models.ModuleAccess{
		RoleID:     roleID,
		ModuleName: req.ModuleName,
		CanView:    req.CanView,
		CanCreate:  req.CanCreate,
		CanEdit:    req.CanEdit,
		CanDelete:  req.CanDelete,
	}

	if err := s.repo.SetModuleAccess(ctx, moduleAccess); err != nil {
		logrus.WithError(err).Error("Failed to set module access")
		return nil, http.StatusInternalServerError, errors.New("failed to set module access")
	}

	return s.mapModuleAccessToResponse(moduleAccess), http.StatusOK, nil
}

func (s *service) GetModuleAccessByRole(ctx context.Context, roleID int) ([]ModuleAccessResponse, int, error) {
	moduleAccess, err := s.repo.GetModuleAccessByRole(ctx, roleID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get module access")
		return nil, http.StatusInternalServerError, errors.New("failed to get module access")
	}

	responses := make([]ModuleAccessResponse, len(moduleAccess))
	for i, ma := range moduleAccess {
		responses[i] = *s.mapModuleAccessToResponse(&ma)
	}

	return responses, http.StatusOK, nil
}

// --- Permission Checking ---

func (s *service) CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, int, error) {
	hasPermission, err := s.repo.UserHasPermission(ctx, userID, permission)
	if err != nil {
		logrus.WithError(err).Error("Failed to check permission")
		return false, http.StatusInternalServerError, errors.New("failed to check permission")
	}

	return hasPermission, http.StatusOK, nil
}

func (s *service) CheckModuleAccess(ctx context.Context, userID uuid.UUID, module string, action string) (bool, int, error) {
	hasAccess, err := s.repo.UserCanAccessModule(ctx, userID, module, action)
	if err != nil {
		logrus.WithError(err).Error("Failed to check module access")
		return false, http.StatusInternalServerError, errors.New("failed to check module access")
	}

	return hasAccess, http.StatusOK, nil
}

// --- Helper Functions ---

func (s *service) mapRoleToResponse(role *models.Role) RoleResponse {
	description := ""
	if role.Description != nil {
		description = *role.Description
	}

	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func (s *service) mapRoleWithPermissionsToResponse(role *models.Role, permissions []models.Permission) *RoleWithPermissionsResponse {
	description := ""
	if role.Description != nil {
		description = *role.Description
	}

	permResponses := make([]PermissionResponse, len(permissions))
	for i, perm := range permissions {
		permResponses[i] = s.mapPermissionToResponse(&perm)
	}

	return &RoleWithPermissionsResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: description,
		IsSystem:    role.IsSystem,
		Permissions: permResponses,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func (s *service) mapPermissionToResponse(permission *models.Permission) PermissionResponse {
	description := ""
	if permission.Description != nil {
		description = *permission.Description
	}

	return PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Module:      permission.Module,
		Description: description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}

func (s *service) mapModuleAccessToResponse(ma *models.ModuleAccess) *ModuleAccessResponse {
	return &ModuleAccessResponse{
		ID:         ma.ID,
		RoleID:     ma.RoleID,
		ModuleName: ma.ModuleName,
		CanView:    ma.CanView,
		CanCreate:  ma.CanCreate,
		CanEdit:    ma.CanEdit,
		CanDelete:  ma.CanDelete,
		CreatedAt:  ma.CreatedAt,
		UpdatedAt:  ma.UpdatedAt,
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
