import { ReactNode } from 'react';
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route';
import { ProtectedModule } from '@/src/presentation/components/layout/protected-feature';

export default function AdminLayout({ children }: { children: ReactNode }) {
    return (
        <ProtectedRoute>
            <ProtectedModule requiredRole={['Super Admin', 'Admin']}>
                {children}
            </ProtectedModule>
        </ProtectedRoute>
    );
}
