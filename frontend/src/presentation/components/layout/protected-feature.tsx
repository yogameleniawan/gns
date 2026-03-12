'use client';

import { usePermission } from '@/src/application/hooks/use-permission';
import { Alert, AlertDescription, AlertTitle } from '@/src/presentation/components/ui/alert';
import { ShieldX } from 'lucide-react';
import { MainLayout } from './main-layout';

interface ProtectedFeatureProps {
    children: React.ReactNode;
    permission?: string;
    role?: string;
    anyRole?: string[];
    requireAdmin?: boolean;
    fallback?: React.ReactNode;
}

/**
 * Wrapper component for feature-level access control
 * Use this to hide/show buttons, tabs, or other UI elements based on permissions
 */
export function ProtectedFeature({
    children,
    permission,
    role,
    anyRole,
    requireAdmin,
    fallback = null,
}: ProtectedFeatureProps) {
    const { hasPermission, hasRole, hasAnyRole, isAdmin } = usePermission();

    // Check admin requirement
    if (requireAdmin && !isAdmin()) {
        return <>{fallback}</>;
    }

    // Check specific permission
    if (permission && !hasPermission(permission)) {
        return <>{fallback}</>;
    }

    // Check specific role
    if (role && !hasRole(role)) {
        return <>{fallback}</>;
    }

    // Check any of the roles
    if (anyRole && !hasAnyRole(anyRole)) {
        return <>{fallback}</>;
    }

    return <>{children}</>;
}

interface ProtectedModuleProps {
    children: React.ReactNode;
    requiredRole?: string | string[];
    requiredPermission?: string;
    showAccessDenied?: boolean;
}

/**
 * Wrapper component for module-level access control
 * Use this to protect entire pages or sections
 */
export function ProtectedModule({
    children,
    requiredRole,
    requiredPermission,
    showAccessDenied = true,
}: ProtectedModuleProps) {
    const { hasPermission, hasRole, hasAnyRole } = usePermission();

    let hasAccess = true;

    if (requiredPermission && !hasPermission(requiredPermission)) {
        hasAccess = false;
    }

    if (requiredRole) {
        if (Array.isArray(requiredRole)) {
            if (!hasAnyRole(requiredRole)) {
                hasAccess = false;
            }
        } else {
            if (!hasRole(requiredRole)) {
                hasAccess = false;
            }
        }
    }

    if (!hasAccess) {
        if (!showAccessDenied) {
            return null;
        }

        return (
            <MainLayout>
                <div className="flex items-center justify-center min-h-[500px] p-8">
                    <div className="relative max-w-md w-full">
                        {/* Card */}
                        <div className="relative bg-card border border-destructive/20 rounded-2xl p-8 shadow-lg backdrop-blur-sm">
                            {/* Icon container */}
                            <div className="flex justify-center mb-6">
                                <div className="relative">
                                    <div className="absolute inset-0 bg-destructive/20 rounded-full blur-xl animate-pulse" />
                                    <div className="relative bg-gradient-to-br from-destructive to-red-600 p-4 rounded-full">
                                        <ShieldX className="h-8 w-8 text-white" />
                                    </div>
                                </div>
                            </div>

                            {/* Content */}
                            <div className="text-center space-y-4">
                                <h2 className="text-2xl font-bold text-foreground">
                                    Access Denied
                                </h2>
                                <p className="text-muted-foreground leading-relaxed">
                                    You don&apos;t have permission to access this module.
                                    Please contact your administrator if you believe this is an error.
                                </p>

                                {/* Decorative divider */}
                                <div className="flex items-center gap-4 py-2">
                                    <div className="flex-1 h-px bg-gradient-to-r from-transparent via-border to-transparent" />
                                    <ShieldX className="h-4 w-4 text-muted-foreground/50" />
                                    <div className="flex-1 h-px bg-gradient-to-r from-transparent via-border to-transparent" />
                                </div>

                                <p className="text-xs text-muted-foreground/70">
                                    Error Code: 403 - Forbidden
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </MainLayout>
        );
    }

    return <>{children}</>;
}
