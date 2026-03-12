'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/src/application/hooks/use-auth';
import { Loader2 } from 'lucide-react';

interface ProtectedRouteProps {
    children: React.ReactNode;
}

/**
 * Wrapper component that requires authentication
 */
export function ProtectedRoute({ children }: ProtectedRouteProps) {
    const { isAuthenticated, isLoading, hasHydrated } = useAuth();
    const router = useRouter();

    useEffect(() => {
        // Only redirect after hydration is complete
        if (hasHydrated && !isLoading && !isAuthenticated) {
            router.push('/auth/login');
        }
    }, [hasHydrated, isAuthenticated, isLoading, router]);

    // Show loading during hydration or while loading
    if (!hasHydrated || isLoading) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        );
    }

    // Redirect in progress
    if (!isAuthenticated) {
        return null;
    }

    return <>{children}</>;
}
