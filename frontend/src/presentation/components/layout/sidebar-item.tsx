'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { ChevronDown, LucideIcon } from 'lucide-react'
import { useState } from 'react'
import { cn } from '@/src/lib/utils'
import { usePermission } from '@/src/application/hooks/use-permission'

export interface SidebarMenuItem {
    label: string
    href?: string
    icon: LucideIcon
    children?: SidebarMenuItem[]
    // RBAC fields - menu item will be hidden if user doesn't have access
    permission?: string           // Specific permission required (e.g., "broker_trades.view")
    role?: string                 // Single role required (e.g., "Super Admin")
    anyRole?: string[]            // Any of these roles required (e.g., ["Super Admin", "Admin"])
    requireAdmin?: boolean        // Requires admin access
}

interface SidebarItemProps {
    item: SidebarMenuItem
    level?: number
    onNavigate?: () => void
}

// Helper function to check if any child matches the current pathname
function hasActiveChild(item: SidebarMenuItem, pathname: string): boolean {
    if (!item.children) return false

    return item.children.some(child => {
        if (child.href === pathname) return true
        if (child.children) return hasActiveChild(child, pathname)
        return false
    })
}

export function SidebarItem({ item, level = 0, onNavigate }: SidebarItemProps) {
    const pathname = usePathname()
    const { hasPermission, hasRole, hasAnyRole, isAdmin } = usePermission()

    const checkAccess = (menuItem: SidebarMenuItem): boolean => {
        // Check admin requirement
        if (menuItem.requireAdmin && !isAdmin()) {
            return false
        }

        // Check specific permission
        if (menuItem.permission && !hasPermission(menuItem.permission)) {
            return false
        }

        // Check specific role
        if (menuItem.role && !hasRole(menuItem.role)) {
            return false
        }

        // Check any of the roles
        if (menuItem.anyRole && !hasAnyRole(menuItem.anyRole)) {
            return false
        }

        return true
    }

    // Filter children based on access
    const accessibleChildren = item.children?.filter(child => checkAccess(child))
    const hasChildren = accessibleChildren && accessibleChildren.length > 0

    // Auto-expand if any child is active
    const shouldBeExpanded = hasChildren && hasActiveChild(item, pathname)
    const [isExpanded, setIsExpanded] = useState(shouldBeExpanded)

    // If user doesn't have access, don't render
    if (!checkAccess(item)) {
        return null
    }

    const isActive = item.href === pathname
    const Icon = item.icon

    const handleClick = () => {
        if (hasChildren) {
            setIsExpanded(!isExpanded)
        } else if (onNavigate) {
            onNavigate()
        }
    }

    const content = (
        <div
            className={cn(
                'flex items-center justify-between gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200',
                'hover:bg-accent hover:text-accent-foreground',
                isActive && 'bg-accent text-accent-foreground shadow-sm',
                level > 0 && 'ml-4'
            )}
        >
            <div className="flex items-center gap-3">
                <Icon className={cn(
                    'h-5 w-5 shrink-0 transition-colors',
                    isActive ? 'text-primary' : 'text-muted-foreground'
                )} />
                <span className="truncate">{item.label}</span>
            </div>
            {hasChildren && (
                <ChevronDown
                    className={cn(
                        'h-4 w-4 shrink-0 text-muted-foreground transition-transform duration-200',
                        isExpanded && 'rotate-180'
                    )}
                />
            )}
        </div>
    )

    return (
        <div>
            {hasChildren ? (
                <button
                    onClick={handleClick}
                    className="w-full text-left"
                >
                    {content}
                </button>
            ) : (
                <Link
                    href={item.href || '#'}
                    onClick={handleClick}
                >
                    {content}
                </Link>
            )}

            {hasChildren && (
                <div
                    className={cn(
                        'overflow-hidden transition-all duration-300',
                        isExpanded ? 'max-h-[1000px] opacity-100' : 'max-h-0 opacity-0'
                    )}
                >
                    <div className="mt-1 space-y-1">
                        {accessibleChildren?.map((child, index) => (
                            <SidebarItem
                                key={index}
                                item={child}
                                level={level + 1}
                                onNavigate={onNavigate}
                            />
                        ))}
                    </div>
                </div>
            )}
        </div>
    )
}
