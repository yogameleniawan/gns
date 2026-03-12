import { useEffect } from 'react';
import { useAuthStore } from '@/src/infrastructure/stores/auth-store';

/**
 * Hook to access authentication state and actions
 */
export function useAuth() {
  const {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
    isLoading,
    hasHydrated,
    login,
    loginWithGoogle,
    register,
    logout,
    refreshAccessToken,
    updateProfile,
    changePassword,
    setLoading,
  } = useAuthStore();

  // Trigger hydration on mount (client-side only)
  useEffect(() => {
    useAuthStore.persist.rehydrate();
  }, []);

  return {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
    isLoading,
    hasHydrated,
    login,
    loginWithGoogle,
    register,
    logout,
    refreshAccessToken,
    updateProfile,
    changePassword,
    setLoading,
  };
}
