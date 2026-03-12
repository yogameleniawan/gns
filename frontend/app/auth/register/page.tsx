'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/src/application/hooks/use-auth';
import { Button } from '@/src/presentation/components/ui/button';
import { Input } from '@/src/presentation/components/ui/input';
import { Label } from '@/src/presentation/components/ui/label';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/src/presentation/components/ui/card';
import { Alert, AlertDescription } from '@/src/presentation/components/ui/alert';
import { Loader2, Check, X } from 'lucide-react';

export default function RegisterPage() {
    const router = useRouter();
    const { register, isLoading } = useAuth();

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const [validations, setValidations] = useState({
        minLength: false,
        hasNumber: false,
        hasSpecial: false,
        match: false,
    });

    const validatePassword = (pwd: string, confirm: string) => {
        setValidations({
            minLength: pwd.length >= 8,
            hasNumber: /\d/.test(pwd),
            hasSpecial: /[!@#$%^&*(),.?":{}|<>]/.test(pwd),
            match: pwd === confirm && pwd.length > 0 && confirm.length > 0,
        });
    };

    const handlePasswordChange = (pwd: string) => {
        setPassword(pwd);
        validatePassword(pwd, confirmPassword);
    };

    const handleConfirmPasswordChange = (confirm: string) => {
        setConfirmPassword(confirm);
        validatePassword(password, confirm);
    };

    const isFormValid = () => {
        return name.trim() !== '' &&
            email.trim() !== '' &&
            validations.minLength &&
            validations.hasNumber &&
            validations.hasSpecial &&
            validations.match;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!isFormValid()) {
            setError('Please ensure all requirements are met');
            return;
        }

        try {
            await register({ name, email, password });
            router.push('/home');
        } catch (err: any) {
            setError(err.response?.data?.message || 'Registration failed. Please try again.');
        }
    };

    const ValidationItem = ({ valid, text }: { valid: boolean; text: string }) => (
        <div className="flex items-center gap-2 text-sm">
            {valid ? (
                <Check className="h-4 w-4 text-green-500" />
            ) : (
                <X className="h-4 w-4 text-gray-300" />
            )}
            <span className={valid ? 'text-green-600' : 'text-gray-500'}>{text}</span>
        </div>
    );

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900 p-4">
            <Card className="w-full max-w-md shadow-xl">
                <CardHeader className="space-y-1">
                    <CardTitle className="text-3xl font-bold text-center">Create Account</CardTitle>
                    <CardDescription className="text-center">
                        Sign up to get started with Quantro
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        {error && (
                            <Alert variant="destructive">
                                <AlertDescription>{error}</AlertDescription>
                            </Alert>
                        )}

                        <div className="space-y-2">
                            <Label htmlFor="name">Full Name</Label>
                            <Input
                                id="name"
                                type="text"
                                placeholder="John Doe"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="email">Email</Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="you@example.com"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="password">Password</Label>
                            <Input
                                id="password"
                                type="password"
                                placeholder="••••••••"
                                value={password}
                                onChange={(e) => handlePasswordChange(e.target.value)}
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="confirmPassword">Confirm Password</Label>
                            <Input
                                id="confirmPassword"
                                type="password"
                                placeholder="••••••••"
                                value={confirmPassword}
                                onChange={(e) => handleConfirmPasswordChange(e.target.value)}
                                required
                                disabled={isLoading}
                            />
                        </div>

                        {/* Password Requirements */}
                        <div className="space-y-2 p-3 bg-muted/50 rounded-md">
                            <p className="text-sm font-medium">Password Requirements:</p>
                            <ValidationItem valid={validations.minLength} text="At least 8 characters" />
                            <ValidationItem valid={validations.hasNumber} text="Contains a number" />
                            <ValidationItem valid={validations.hasSpecial} text="Contains a special character" />
                            <ValidationItem valid={validations.match} text="Passwords match" />
                        </div>

                        <Button
                            type="submit"
                            className="w-full"
                            disabled={isLoading || !isFormValid()}
                        >
                            {isLoading ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Creating account...
                                </>
                            ) : (
                                'Create Account'
                            )}
                        </Button>
                    </form>
                </CardContent>
                <CardFooter>
                    <p className="text-center text-sm text-muted-foreground w-full">
                        Already have an account?{' '}
                        <Link href="/auth/login" className="text-primary hover:underline font-medium">
                            Sign in
                        </Link>
                    </p>
                </CardFooter>
            </Card>
        </div>
    );
}
