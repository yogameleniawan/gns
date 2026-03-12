import apiClient from './api-client';

/**
 * Role interface
 */
export interface Role {
  id: number;
  name: string;
  description?: string;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}

/**
 * Permission interface
 */
export interface Permission {
  id: number;
  name: string;
  module: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

/**
 * Role with permissions
 */
export interface RoleWithPermissions extends Role {
  permissions: Permission[];
}

/**
 * Permissions grouped by module
 */
export interface PermissionsByModule {
  module: string;
  permissions: Permission[];
}

/**
 * Module access interface
 */
export interface ModuleAccess {
  id: number;
  role_id: number;
  module_name: string;
  can_view: boolean;
  can_create: boolean;
  can_edit: boolean;
  can_delete: boolean;
  created_at: string;
  updated_at: string;
}

/**
 * Create role request
 */
export interface CreateRoleRequest {
  name: string;
  description?: string;
}

/**
 * Update role request
 */
export interface UpdateRoleRequest {
  name?: string;
  description?: string;
}

/**
 * Create permission request
 */
export interface CreatePermissionRequest {
  name: string;
  module: string;
  description?: string;
}

/**
 * Update permission request
 */
export interface UpdatePermissionRequest {
  name?: string;
  module?: string;
  description?: string;
}

/**
 * Assign permissions to role request
 */
export interface AssignPermissionsRequest {
  permission_ids: number[];
}

/**
 * Assign roles to user request
 */
export interface AssignRolesRequest {
  role_ids: number[];
}

/**
 * Update module access request
 */
export interface UpdateModuleAccessRequest {
  module_name: string;
  can_view: boolean;
  can_create: boolean;
  can_edit: boolean;
  can_delete: boolean;
}

/**
 * RBAC Service
 */
export const rbacService = {
  // ========== Roles ==========
  
  /**
   * Get all roles
   */
  async getAllRoles(): Promise<Role[]> {
    const response = await apiClient.get<Role[]>('/v1/rbac/roles');
    return response.data;
  },

  /**
   * Get role by ID with permissions
   */
  async getRoleById(id: number): Promise<RoleWithPermissions> {
    const response = await apiClient.get<RoleWithPermissions>(`/v1/rbac/roles/${id}`);
    return response.data;
  },

  /**
   * Create a new role
   */
  async createRole(data: CreateRoleRequest): Promise<Role> {
    const response = await apiClient.post<Role>('/v1/rbac/roles', data);
    return response.data;
  },

  /**
   * Update role
   */
  async updateRole(id: number, data: UpdateRoleRequest): Promise<Role> {
    const response = await apiClient.put<Role>(`/v1/rbac/roles/${id}`, data);
    return response.data;
  },

  /**
   * Delete role
   */
  async deleteRole(id: number): Promise<void> {
    await apiClient.delete(`/v1/rbac/roles/${id}`);
  },

  // ========== Permissions ==========

  /**
   * Get all permissions
   */
  async getAllPermissions(): Promise<Permission[]> {
    const response = await apiClient.get<Permission[]>('/v1/rbac/permissions');
    return response.data;
  },

  /**
   * Get permissions grouped by module
   */
  async getPermissionsByModule(): Promise<PermissionsByModule[]> {
    const response = await apiClient.get<PermissionsByModule[]>('/v1/rbac/permissions/by-module');
    return response.data;
  },

  /**
   * Get permission by ID
   */
  async getPermissionById(id: number): Promise<Permission> {
    const response = await apiClient.get<Permission>(`/v1/rbac/permissions/${id}`);
    return response.data;
  },

  /**
   * Create a new permission
   */
  async createPermission(data: CreatePermissionRequest): Promise<Permission> {
    const response = await apiClient.post<Permission>('/v1/rbac/permissions', data);
    return response.data;
  },

  /**
   * Update permission
   */
  async updatePermission(id: number, data: UpdatePermissionRequest): Promise<Permission> {
    const response = await apiClient.put<Permission>(`/v1/rbac/permissions/${id}`, data);
    return response.data;
  },

  /**
   * Delete permission
   */
  async deletePermission(id: number): Promise<void> {
    await apiClient.delete(`/v1/rbac/permissions/${id}`);
  },

  // ========== Role-Permission Assignment ==========

  /**
   * Assign permissions to role
   */
  async assignPermissionsToRole(roleId: number, data: AssignPermissionsRequest): Promise<void> {
    await apiClient.post(`/v1/rbac/roles/${roleId}/permissions`, data);
  },

  /**
   * Get role permissions
   */
  async getRolePermissions(roleId: number): Promise<Permission[]> {
    const response = await apiClient.get<Permission[]>(`/v1/rbac/roles/${roleId}/permissions`);
    return response.data;
  },

  // ========== User-Role Assignment ==========

  /**
   * Assign roles to user
   */
  async assignRolesToUser(userId: string, data: AssignRolesRequest): Promise<void> {
    await apiClient.post(`/v1/rbac/users/${userId}/roles`, data);
  },

  /**
   * Get user roles
   */
  async getUserRoles(userId: string): Promise<Role[]> {
    const response = await apiClient.get<Role[]>(`/v1/rbac/users/${userId}/roles`);
    return response.data;
  },

  // ========== Module Access ==========

  /**
   * Update module access for role
   */
  async updateModuleAccess(roleId: number, data: UpdateModuleAccessRequest): Promise<ModuleAccess> {
    const response = await apiClient.post<ModuleAccess>(`/v1/rbac/roles/${roleId}/module-access`, data);
    return response.data;
  },

  /**
   * Get module access by role
   */
  async getModuleAccessByRole(roleId: number): Promise<ModuleAccess[]> {
    const response = await apiClient.get<ModuleAccess[]>(`/v1/rbac/roles/${roleId}/module-access`);
    return response.data;
  },

  // ========== Permission Checking ==========

  /**
   * Check if current user has permission
   */
  async checkPermission(permission: string): Promise<boolean> {
    const response = await apiClient.post<{ has_permission: boolean }>('/v1/rbac/check-permission', {
      permission,
    });
    return response.data.has_permission;
  },

  /**
   * Check if current user can access module
   */
  async checkModuleAccess(module: string, action: 'view' | 'create' | 'edit' | 'delete'): Promise<boolean> {
    const response = await apiClient.post<{ has_access: boolean }>('/v1/rbac/check-module-access', {
      module,
      action,
    });
    return response.data.has_access;
  },
};
