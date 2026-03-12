package rbac

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/database"
)

var (
	ErrRoleNotFound        = errors.New("role not found")
	ErrRoleAlreadyExists   = errors.New("role already exists")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrSystemRoleProtected = errors.New("cannot modify system role")
)

type Repository interface {
	// Role operations
	CreateRole(ctx context.Context, role *models.Role) error
	GetRoleByID(ctx context.Context, id int) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	GetAllRoles(ctx context.Context) ([]models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id int) error

	// Permission operations
	CreatePermission(ctx context.Context, permission *models.Permission) error
	GetPermissionByID(ctx context.Context, id int) (*models.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*models.Permission, error)
	GetAllPermissions(ctx context.Context) ([]models.Permission, error)
	GetPermissionsByModule(ctx context.Context, module string) ([]models.Permission, error)
	GetPermissionsByModuleGrouped(ctx context.Context) (map[string][]models.Permission, error)
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	DeletePermission(ctx context.Context, id int) error

	// Role-Permission operations
	AssignPermissionsToRole(ctx context.Context, roleID int, permissionIDs []int) error
	RevokePermissionsFromRole(ctx context.Context, roleID int, permissionIDs []int) error
	GetRolePermissions(ctx context.Context, roleID int) ([]models.Permission, error)

	// User-Role operations
	AssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []int) error
	RevokeRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []int) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error)

	// Module access operations
	SetModuleAccess(ctx context.Context, moduleAccess *models.ModuleAccess) error
	GetModuleAccessByRole(ctx context.Context, roleID int) ([]models.ModuleAccess, error)
	GetModuleAccessByRoleAndModule(ctx context.Context, roleID int, moduleName string) (*models.ModuleAccess, error)

	// Permission checking
	UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	UserHasRole(ctx context.Context, userID uuid.UUID, role string) (bool, error)
	UserCanAccessModule(ctx context.Context, userID uuid.UUID, module string, action string) (bool, error)
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{
		db: db,
	}
}

// --- Role Operations ---

func (r *repository) CreateRole(ctx context.Context, role *models.Role) error {
	var existing models.Role
	if err := r.db.GetDB().Where("name = ?", role.Name).First(&existing).Error; err == nil {
		return ErrRoleAlreadyExists
	}

	if err := r.db.GetDB().WithContext(ctx).Create(role).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetRoleByID(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role
	if err := r.db.GetDB().WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.GetDB().WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) UpdateRole(ctx context.Context, role *models.Role) error {
	if err := r.db.GetDB().WithContext(ctx).Save(role).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteRole(ctx context.Context, id int) error {
	// Check if it's a system role
	role, err := r.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return ErrSystemRoleProtected
	}

	if err := r.db.GetDB().WithContext(ctx).Delete(&models.Role{}, id).Error; err != nil {
		return err
	}
	return nil
}

// --- Permission Operations ---

func (r *repository) CreatePermission(ctx context.Context, permission *models.Permission) error {
	if err := r.db.GetDB().WithContext(ctx).Create(permission).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetPermissionByID(ctx context.Context, id int) (*models.Permission, error) {
	var permission models.Permission
	if err := r.db.GetDB().WithContext(ctx).Where("id = ?", id).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

func (r *repository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	var permission models.Permission
	if err := r.db.GetDB().WithContext(ctx).Where("name = ?", name).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

func (r *repository) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.GetDB().WithContext(ctx).Order("module, name").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *repository) GetPermissionsByModule(ctx context.Context, module string) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.GetDB().WithContext(ctx).
		Where("module = ?", module).
		Order("name").
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *repository) GetPermissionsByModuleGrouped(ctx context.Context) (map[string][]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.GetDB().WithContext(ctx).Order("module, name").Find(&permissions).Error; err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.Permission)
	for _, perm := range permissions {
		grouped[perm.Module] = append(grouped[perm.Module], perm)
	}

	return grouped, nil
}

func (r *repository) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	if err := r.db.GetDB().WithContext(ctx).Save(permission).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) DeletePermission(ctx context.Context, id int) error {
	if err := r.db.GetDB().WithContext(ctx).Delete(&models.Permission{}, id).Error; err != nil {
		return err
	}
	return nil
}

// --- Role-Permission Operations ---

func (r *repository) AssignPermissionsToRole(ctx context.Context, roleID int, permissionIDs []int) error {
	// First, delete existing assignments
	if err := r.db.GetDB().WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	// Then create new assignments
	for _, permID := range permissionIDs {
		rolePermission := &models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := r.db.GetDB().WithContext(ctx).Create(rolePermission).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) RevokePermissionsFromRole(ctx context.Context, roleID int, permissionIDs []int) error {
	if err := r.db.GetDB().WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetRolePermissions(ctx context.Context, roleID int) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := r.db.GetDB().WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.module, permissions.name").
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// --- User-Role Operations ---

func (r *repository) AssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []int) error {
	// First, delete existing assignments
	if err := r.db.GetDB().WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.UserRole{}).Error; err != nil {
		return err
	}

	// Then create new assignments
	for _, roleID := range roleIDs {
		userRole := &models.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := r.db.GetDB().WithContext(ctx).Create(userRole).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) RevokeRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []int) error {
	if err := r.db.GetDB().WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&models.UserRole{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.GetDB().WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// --- Module Access Operations ---

func (r *repository) SetModuleAccess(ctx context.Context, moduleAccess *models.ModuleAccess) error {
	// Upsert: Update if exists, insert if not
	var existing models.ModuleAccess
	err := r.db.GetDB().WithContext(ctx).
		Where("role_id = ? AND module_name = ?", moduleAccess.RoleID, moduleAccess.ModuleName).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Insert
		if err := r.db.GetDB().WithContext(ctx).Create(moduleAccess).Error; err != nil {
			return err
		}
	} else {
		// Update
		existing.CanView = moduleAccess.CanView
		existing.CanCreate = moduleAccess.CanCreate
		existing.CanEdit = moduleAccess.CanEdit
		existing.CanDelete = moduleAccess.CanDelete
		if err := r.db.GetDB().WithContext(ctx).Save(&existing).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) GetModuleAccessByRole(ctx context.Context, roleID int) ([]models.ModuleAccess, error) {
	var moduleAccess []models.ModuleAccess
	if err := r.db.GetDB().WithContext(ctx).
		Where("role_id = ?", roleID).
		Order("module_name").
		Find(&moduleAccess).Error; err != nil {
		return nil, err
	}
	return moduleAccess, nil
}

func (r *repository) GetModuleAccessByRoleAndModule(ctx context.Context, roleID int, moduleName string) (*models.ModuleAccess, error) {
	var moduleAccess models.ModuleAccess
	if err := r.db.GetDB().WithContext(ctx).
		Where("role_id = ? AND module_name = ?", roleID, moduleName).
		First(&moduleAccess).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &moduleAccess, nil
}

// --- Permission Checking ---

func (r *repository) UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	var count int64
	if err := r.db.GetDB().WithContext(ctx).
		Table("user_roles").
		Joins("JOIN role_permissions ON role_permissions.role_id = user_roles.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ? AND permissions.name = ?", userID, permission).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) UserHasRole(ctx context.Context, userID uuid.UUID, role string) (bool, error) {
	var count int64
	if err := r.db.GetDB().WithContext(ctx).
		Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, role).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) UserCanAccessModule(ctx context.Context, userID uuid.UUID, module string, action string) (bool, error) {
	query := r.db.GetDB().WithContext(ctx).
		Table("user_roles").
		Joins("JOIN module_access ON module_access.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND module_access.module_name = ?", userID, module)

	// Add action-specific condition
	switch action {
	case "view":
		query = query.Where("module_access.can_view = ?", true)
	case "create":
		query = query.Where("module_access.can_create = ?", true)
	case "edit":
		query = query.Where("module_access.can_edit = ?", true)
	case "delete":
		query = query.Where("module_access.can_delete = ?", true)
	default:
		return false, errors.New("invalid action")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
