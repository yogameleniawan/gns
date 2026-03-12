import { useAuthStore } from '@/src/infrastructure/stores/auth-store';

/**
 * Hook for permission checking
 */
export function usePermission() {
  const { user } = useAuthStore();

  /**
   * Check if user has a specific permission
   */
  const hasPermission = (permission: string): boolean => {
    if (!user || !user.permissions) return false;
    return user.permissions.includes(permission);
  };

  /**
   * Check if user has any of the specified permissions
   */
  const hasAnyPermission = (permissions: string[]): boolean => {
    if (!user || !user.permissions) return false;
    return permissions.some(p => user.permissions?.includes(p));
  };

  /**
   * Check if user has all specified permissions
   */
  const hasAllPermissions = (permissions: string[]): boolean => {
    if (!user || !user.permissions) return false;
    return permissions.every(p => user.permissions?.includes(p));
  };

  /**
   * Check if user has a specific role
   */
  const hasRole = (role: string): boolean => {
    if (!user || !user.roles) return false;
    return user.roles.includes(role);
  };

  /**
   * Check if user has any of the specified roles
   */
  const hasAnyRole = (roles: string[]): boolean => {
    if (!user || !user.roles) return false;
    return roles.some(r => user.roles?.includes(r));
  };

  /**
   * Check if user is admin (has "Super Admin" or "Admin" role)
   */
  const isAdmin = (): boolean => {
    return hasAnyRole(['Super Admin', 'Admin']);
  };

  /**
   * Check if user is super admin
   */
  const isSuperAdmin = (): boolean => {
    return hasRole('Super Admin');
  };

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
    hasRole,
    hasAnyRole,
    isAdmin,
    isSuperAdmin,
  };
}
