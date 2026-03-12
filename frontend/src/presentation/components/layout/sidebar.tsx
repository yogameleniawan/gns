'use client'

import { useState } from 'react'
import {
    X, LayoutDashboard, Calculator, DollarSign, TrendingUp, TrendingDown, Briefcase, User, Settings,
    Shield, Target, Repeat, Receipt, PieChart, GitBranch, Scissors, Activity, ArrowUpCircle, LineChart, LogOut,
    Users, ShieldCheck, Key, AlertTriangle
} from 'lucide-react'
import { SidebarItem, SidebarMenuItem } from './sidebar-item'
import { cn } from '@/src/lib/utils'
import { Button } from '@/src/presentation/components/ui/button'
import { useAuth } from '@/src/application/hooks/use-auth'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/src/presentation/components/ui/dropdown-menu'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/src/presentation/components/ui/dialog'

interface SidebarProps {
    isOpen: boolean
    onClose: () => void
}

const menuItems: SidebarMenuItem[] = [
    {
        label: 'Home',
        href: '/home',
        icon: LayoutDashboard,
    },
    {
        label: 'Broker Flow',
        href: '/broker-flow',
        icon: Briefcase,
    },
    {
        label: 'Stock Flow',
        href: '/stock-flow',
        icon: Activity,
    },
    {
        label: 'Sentiment Indicator',
        href: '/sentiment-indicator',
        icon: TrendingUp,
    },
    {
        label: 'Market',
        href: '/market',
        icon: LineChart,
    },
    {
        label: 'Calculator',
        icon: Calculator,
        children: [
            {
                label: 'Profit Calculator',
                href: '/calculator/profit',
                icon: DollarSign,
            },
            {
                label: 'Dividend Calculator',
                href: '/calculator/dividend',
                icon: TrendingUp,
            },
            {
                label: 'ARA/ARB Calculator',
                href: '/calculator/ara-arb',
                icon: Activity,
            },
            {
                label: 'Lot & Fee Calculator',
                href: '/calculator/lot-fee',
                icon: Receipt,
            },
            {
                label: 'Position Sizing',
                href: '/calculator/position-sizing',
                icon: Shield,
            },
            {
                label: 'Risk Reward',
                href: '/calculator/risk-reward',
                icon: Target,
            },
            {
                label: 'Average Down/Up',
                href: '/calculator/average',
                icon: TrendingDown,
            },
            {
                label: 'DCA Calculator',
                href: '/calculator/dca',
                icon: Repeat,
            },
            {
                label: 'Portfolio Diversification',
                href: '/calculator/portfolio',
                icon: PieChart,
            },
            {
                label: 'Breakeven Calculator',
                href: '/calculator/breakeven',
                icon: ArrowUpCircle,
            },
            {
                label: 'Rights Issue (HMETD)',
                href: '/calculator/rights-issue',
                icon: GitBranch,
            },
            {
                label: 'Stock Split',
                href: '/calculator/stock-split',
                icon: Scissors,
            },
            {
                label: 'Moving Average',
                href: '/calculator/moving-average',
                icon: Activity,
            },
            {
                label: 'Investment Growth',
                href: '/calculator/investment-growth',
                icon: TrendingUp,
            },
        ],
    },
    {
        label: 'User Management',
        icon: Shield,
        anyRole: ['Super Admin', 'Admin'],
        children: [
            {
                label: 'Users',
                href: '/admin/users',
                icon: Users,
                permission: 'users.read',
            },
            {
                label: 'Roles',
                href: '/admin/roles',
                icon: ShieldCheck,
                permission: 'roles.read',
            },
            {
                label: 'Permissions',
                href: '/admin/permissions',
                icon: Key,
                permission: 'permissions.read',
            },
        ],
    },
]

export function Sidebar({ isOpen, onClose }: SidebarProps) {
    const { user, logout } = useAuth()
    const [isLogoutModalOpen, setIsLogoutModalOpen] = useState(false)

    // Get user initials for avatar
    const getUserInitials = () => {
        if (!user?.name) return 'U'
        const names = user.name.split(' ')
        if (names.length >= 2) {
            return `${names[0][0]}${names[1][0]}`.toUpperCase()
        }
        return user.name.substring(0, 2).toUpperCase()
    }

    const handleLogout = () => {
        setIsLogoutModalOpen(true)
    }

    const confirmLogout = () => {
        logout()
        setIsLogoutModalOpen(false)
    }

    return (
        <>
            {/* Backdrop for mobile/tablet */}
            <div
                className={cn(
                    'fixed inset-0 z-40 bg-black/50 backdrop-blur-sm transition-opacity duration-300 lg:hidden',
                    isOpen ? 'opacity-100' : 'pointer-events-none opacity-0'
                )}
                onClick={onClose}
            />

            {/* Sidebar */}
            <aside
                className={cn(
                    'fixed left-0 top-0 z-50 flex h-screen w-[280px] flex-col border-r bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60 transition-transform duration-300 ease-in-out',
                    // Mobile: slide in from left
                    'lg:translate-x-0',
                    isOpen ? 'translate-x-0' : '-translate-x-full'
                )}
            >
                {/* Sidebar Header */}
                <div className="flex h-16 items-center justify-between border-b px-4">
                    <div className="flex items-center space-x-2">
                        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                            <span className="text-xl font-bold">Q</span>
                        </div>
                        <span className="text-xl font-bold">Quantro</span>
                    </div>

                    {/* Close button (mobile/tablet only) */}
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={onClose}
                        className="lg:hidden"
                        aria-label="Close sidebar"
                    >
                        <X className="h-5 w-5" />
                    </Button>
                </div>

                {/* Navigation */}
                <nav className="flex-1 space-y-1 overflow-y-auto p-4">
                    {menuItems.map((item, index) => (
                        <SidebarItem
                            key={index}
                            item={item}
                            onNavigate={onClose}
                        />
                    ))}
                </nav>

                {/* Footer - User Profile */}
                <div className="border-t p-4">
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <button className="flex w-full items-center gap-3 rounded-lg bg-muted/50 px-3 py-2 transition-colors hover:bg-muted">
                                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground text-sm font-semibold">
                                    {getUserInitials()}
                                </div>
                                <div className="flex-1 overflow-hidden text-left">
                                    <p className="truncate text-sm font-medium">
                                        {user?.name || 'User Name'}
                                    </p>
                                    <p className="truncate text-xs text-muted-foreground">
                                        {user?.email || 'user@example.com'}
                                    </p>
                                </div>
                            </button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-56">
                            <DropdownMenuItem asChild>
                                <a href="/settings/profile" className="cursor-pointer">
                                    <User className="mr-2 h-4 w-4" />
                                    <span>Profile</span>
                                </a>
                            </DropdownMenuItem>
                            <DropdownMenuItem asChild>
                                <a href="/settings/preference" className="cursor-pointer">
                                    <Settings className="mr-2 h-4 w-4" />
                                    <span>Settings</span>
                                </a>
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                                onClick={handleLogout}
                                className="cursor-pointer text-red-600 focus:text-red-600"
                            >
                                <LogOut className="mr-2 h-4 w-4" />
                                <span>Logout</span>
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </aside>

            {/* Logout Confirmation Modal */}
            <Dialog open={isLogoutModalOpen} onOpenChange={setIsLogoutModalOpen}>
                <DialogContent className="sm:max-w-[425px]">
                    <DialogHeader>
                        <div className="flex items-center justify-center w-12 h-12 mx-auto mb-4 rounded-full bg-red-100 dark:bg-red-900/20">
                            <AlertTriangle className="w-6 h-6 text-red-600 dark:text-red-500" />
                        </div>
                        <DialogTitle className="text-center text-xl">
                            Konfirmasi Logout
                        </DialogTitle>
                        <DialogDescription className="text-center">
                            Apakah Anda yakin ingin keluar dari aplikasi?
                            <br />
                            <span className="text-sm text-muted-foreground mt-2 block">
                                Anda perlu login kembali untuk mengakses aplikasi.
                            </span>
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter className="flex-col sm:flex-row gap-2 sm:gap-0">
                        <Button
                            variant="outline"
                            onClick={() => setIsLogoutModalOpen(false)}
                            className="w-full sm:w-auto"
                        >
                            Batal
                        </Button>
                        <Button
                            variant="destructive"
                            onClick={confirmLogout}
                            className="w-full sm:w-auto"
                        >
                            <LogOut className="mr-2 h-4 w-4" />
                            Ya, Logout
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    )
}
