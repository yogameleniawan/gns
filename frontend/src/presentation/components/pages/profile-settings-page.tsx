'use client'

import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { ProfileInfoForm } from '@/src/presentation/components/forms/profile-info-form'
import { ChangePasswordForm } from '@/src/presentation/components/forms/change-password-form'
import { useAuth } from '@/src/application/hooks/use-auth'
import { Card, CardContent, CardHeader } from '@/src/presentation/components/ui/card'
import { Separator } from '@/src/presentation/components/ui/separator'

export function ProfileSettingsPage() {
    const { user } = useAuth()

    // Get user initials for avatar
    const getUserInitials = () => {
        if (!user?.name) return 'U'
        const names = user.name.split(' ')
        if (names.length >= 2) {
            return `${names[0][0]}${names[1][0]}`.toUpperCase()
        }
        return user.name.substring(0, 2).toUpperCase()
    }

    return (
        <ProtectedRoute>
            <MainLayout>
                <div className="space-y-6 pb-8">
                    {/* Page Header */}
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">Profile Settings</h1>
                        <p className="text-muted-foreground mt-2">
                            Manage your account information and security settings
                        </p>
                    </div>

                    <Separator />

                    {/* User Info Card */}
                    <Card>
                        <CardHeader>
                            <div className="flex items-center gap-4">
                                <div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-primary-foreground text-xl font-semibold">
                                    {getUserInitials()}
                                </div>
                                <div>
                                    <h2 className="text-xl font-semibold">{user?.name || 'User Name'}</h2>
                                    <p className="text-sm text-muted-foreground">{user?.email || 'user@example.com'}</p>
                                    {user?.is_oauth && user?.oauth_provider && (
                                        <p className="text-xs text-muted-foreground mt-1">
                                            Signed in with {user.oauth_provider}
                                        </p>
                                    )}
                                </div>
                            </div>
                        </CardHeader>
                    </Card>

                    {/* Forms Grid */}
                    <div className="grid gap-6 lg:grid-cols-2">
                        <ProfileInfoForm />

                        {/* Only show password change for non-OAuth users */}
                        {!user?.is_oauth && <ChangePasswordForm />}
                    </div>

                    {/* OAuth Notice */}
                    {user?.is_oauth && (
                        <Card className="border-muted">
                            <CardContent className="pt-6">
                                <p className="text-sm text-muted-foreground">
                                    You are signed in with {user.oauth_provider}. Password management is handled by your OAuth provider.
                                </p>
                            </CardContent>
                        </Card>
                    )}
                </div>
            </MainLayout>
        </ProtectedRoute>
    )
}
