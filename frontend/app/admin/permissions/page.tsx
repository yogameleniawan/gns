import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { ProtectedModule } from '@/src/presentation/components/layout/protected-feature'
import { PermissionsContainer } from '@/src/presentation/components/admin/permissions-container'

export default function PermissionsPage() {
    return (
        <ProtectedRoute>
            <ProtectedModule requiredRole={['Super Admin', 'Admin']}>
                <MainLayout>
                    <PermissionsContainer />
                </MainLayout>
            </ProtectedModule>
        </ProtectedRoute>
    )
}
