'use client'

import { useState } from 'react'
import { Button } from '@/src/presentation/components/ui/button'
import { Input } from '@/src/presentation/components/ui/input'
import { Label } from '@/src/presentation/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/src/presentation/components/ui/card'
import { useAuth } from '@/src/application/hooks/use-auth'
import { Loader2, User } from 'lucide-react'
import { toast } from 'sonner'

export function ProfileInfoForm() {
    const { user, updateProfile, isLoading } = useAuth()
    const [name, setName] = useState(user?.name || '')
    const [isSubmitting, setIsSubmitting] = useState(false)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        // Validation
        if (!name.trim()) {
            toast.error('Name is required')
            return
        }

        if (name.trim().length < 2) {
            toast.error('Name must be at least 2 characters')
            return
        }

        if (name === user?.name) {
            toast.info('No changes detected')
            return
        }

        try {
            setIsSubmitting(true)
            await updateProfile({ name: name.trim() })
            toast.success('Profile updated successfully!')
        } catch (error: any) {
            toast.error(error.response?.data?.message || 'Failed to update profile')
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <User className="h-5 w-5" />
                    Profile Information
                </CardTitle>
                <CardDescription>
                    Update your personal information
                </CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="email">Email</Label>
                        <Input
                            id="email"
                            type="email"
                            value={user?.email || ''}
                            disabled
                            className="bg-muted cursor-not-allowed"
                        />
                        <p className="text-xs text-muted-foreground">
                            Email cannot be changed
                        </p>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="name">Name *</Label>
                        <Input
                            id="name"
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Your full name"
                            required
                            minLength={2}
                            disabled={isSubmitting}
                        />
                    </div>

                    <Button
                        type="submit"
                        disabled={isSubmitting || isLoading}
                        className="w-full sm:w-auto"
                    >
                        {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        Update Profile
                    </Button>
                </form>
            </CardContent>
        </Card>
    )
}
