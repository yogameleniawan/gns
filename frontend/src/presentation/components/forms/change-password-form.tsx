'use client'

import { useState } from 'react'
import { Button } from '@/src/presentation/components/ui/button'
import { Input } from '@/src/presentation/components/ui/input'
import { Label } from '@/src/presentation/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/src/presentation/components/ui/card'
import { useAuth } from '@/src/application/hooks/use-auth'
import { Loader2, Lock, Eye, EyeOff } from 'lucide-react'
import { toast } from 'sonner'

export function ChangePasswordForm() {
    const { changePassword, isLoading } = useAuth()
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [showOldPassword, setShowOldPassword] = useState(false)
    const [showNewPassword, setShowNewPassword] = useState(false)
    const [showConfirmPassword, setShowConfirmPassword] = useState(false)

    const [formData, setFormData] = useState({
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
    })

    const handleChange = (field: keyof typeof formData) => (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData(prev => ({ ...prev, [field]: e.target.value }))
    }

    const resetForm = () => {
        setFormData({
            oldPassword: '',
            newPassword: '',
            confirmPassword: '',
        })
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        // Validation
        if (!formData.oldPassword || !formData.newPassword || !formData.confirmPassword) {
            toast.error('All fields are required')
            return
        }

        if (formData.newPassword.length < 8) {
            toast.error('New password must be at least 8 characters')
            return
        }

        if (formData.newPassword === formData.oldPassword) {
            toast.error('New password must be different from old password')
            return
        }

        if (formData.newPassword !== formData.confirmPassword) {
            toast.error('New password and confirmation do not match')
            return
        }

        try {
            setIsSubmitting(true)
            await changePassword({
                old_password: formData.oldPassword,
                new_password: formData.newPassword,
            })
            toast.success('Password changed successfully! You will be logged out.')
            resetForm()
            // Note: User will be automatically logged out by the auth store
        } catch (error: any) {
            toast.error(error.response?.data?.message || 'Failed to change password')
            setIsSubmitting(false)
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Lock className="h-5 w-5" />
                    Change Password
                </CardTitle>
                <CardDescription>
                    Update your password to keep your account secure
                </CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="oldPassword">Current Password *</Label>
                        <div className="relative">
                            <Input
                                id="oldPassword"
                                type={showOldPassword ? 'text' : 'password'}
                                value={formData.oldPassword}
                                onChange={handleChange('oldPassword')}
                                placeholder="Enter current password"
                                required
                                disabled={isSubmitting}
                                className="pr-10"
                            />
                            <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                onClick={() => setShowOldPassword(!showOldPassword)}
                                tabIndex={-1}
                            >
                                {showOldPassword ? (
                                    <EyeOff className="h-4 w-4 text-muted-foreground" />
                                ) : (
                                    <Eye className="h-4 w-4 text-muted-foreground" />
                                )}
                            </Button>
                        </div>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="newPassword">New Password *</Label>
                        <div className="relative">
                            <Input
                                id="newPassword"
                                type={showNewPassword ? 'text' : 'password'}
                                value={formData.newPassword}
                                onChange={handleChange('newPassword')}
                                placeholder="Enter new password (min. 8 characters)"
                                required
                                minLength={8}
                                disabled={isSubmitting}
                                className="pr-10"
                            />
                            <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                onClick={() => setShowNewPassword(!showNewPassword)}
                                tabIndex={-1}
                            >
                                {showNewPassword ? (
                                    <EyeOff className="h-4 w-4 text-muted-foreground" />
                                ) : (
                                    <Eye className="h-4 w-4 text-muted-foreground" />
                                )}
                            </Button>
                        </div>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="confirmPassword">Confirm New Password *</Label>
                        <div className="relative">
                            <Input
                                id="confirmPassword"
                                type={showConfirmPassword ? 'text' : 'password'}
                                value={formData.confirmPassword}
                                onChange={handleChange('confirmPassword')}
                                placeholder="Confirm new password"
                                required
                                disabled={isSubmitting}
                                className="pr-10"
                            />
                            <Button
                                type="button"
                                variant="ghost"
                                size="sm"
                                className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                tabIndex={-1}
                            >
                                {showConfirmPassword ? (
                                    <EyeOff className="h-4 w-4 text-muted-foreground" />
                                ) : (
                                    <Eye className="h-4 w-4 text-muted-foreground" />
                                )}
                            </Button>
                        </div>
                    </div>

                    <div className="pt-2">
                        <Button
                            type="submit"
                            variant="destructive"
                            disabled={isSubmitting || isLoading}
                            className="w-full sm:w-auto"
                        >
                            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            Change Password
                        </Button>
                        <p className="text-xs text-muted-foreground mt-2">
                            You will be logged out after changing your password
                        </p>
                    </div>
                </form>
            </CardContent>
        </Card>
    )
}
