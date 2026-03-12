import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authService, type User, type LoginRequest, type RegisterRequest, type GoogleOAuthRequest, type UpdateProfileRequest, type ChangePasswordRequest } from '@/src/domain/services/auth.service';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  hasHydrated: boolean;
  
  // Actions
  login: (credentials: LoginRequest) => Promise<void>;
  loginWithGoogle: (data: GoogleOAuthRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
  updateProfile: (data: UpdateProfileRequest) => Promise<void>;
  changePassword: (data: ChangePasswordRequest) => Promise<void>;
  setUser: (user: User | null) => void;
  setTokens: (accessToken: string, refreshToken: string) => void;
  setLoading: (loading: boolean) => void;
  setHasHydrated: (state: boolean) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
      hasHydrated: false,

      login: async (credentials: LoginRequest) => {
        try {
          set({ isLoading: true });
          const response = await authService.login(credentials);
          
          set({
            user: response.user,
            accessToken: response.access_token,
            refreshToken: response.refresh_token,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      loginWithGoogle: async (data: GoogleOAuthRequest) => {
        try {
          set({ isLoading: true });
          const response = await authService.loginWithGoogle(data);
          
          set({
            user: response.user,
            accessToken: response.access_token,
            refreshToken: response.refresh_token,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      register: async (data: RegisterRequest) => {
        try {
          set({ isLoading: true });
          const response = await authService.register(data);
          
          set({
            user: response.user,
            accessToken: response.access_token,
            refreshToken: response.refresh_token,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      logout: () => {
        const { accessToken } = get();
        
        // Call logout API if user is authenticated
        if (accessToken) {
          authService.logout().catch(console.error);
        }
        
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        });
        
        // Redirect to login
        if (typeof window !== 'undefined') {
          window.location.href = '/auth/login';
        }
      },

      refreshAccessToken: async () => {
        try {
          const { refreshToken } = get();
          
          if (!refreshToken) {
            throw new Error('No refresh token available');
          }
          
          const response = await authService.refreshToken({ refresh_token: refreshToken });
          
          set({
            user: response.user,
            accessToken: response.access_token,
            refreshToken: response.refresh_token,
            isAuthenticated: true,
          });
        } catch (error) {
          // Refresh failed, logout
          get().logout();
          throw error;
        }
      },

      updateProfile: async (data: UpdateProfileRequest) => {
        try {
          set({ isLoading: true });
          const updatedUser = await authService.updateProfile(data);
          
          set({
            user: updatedUser,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      changePassword: async (data: ChangePasswordRequest) => {
        try {
          set({ isLoading: true });
          await authService.changePassword(data);
          set({ isLoading: false });
          
          // Logout after password change (user will need to login again)
          get().logout();
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      setUser: (user: User | null) => {
        set({ user, isAuthenticated: !!user });
      },

      setTokens: (accessToken: string, refreshToken: string) => {
        set({ accessToken, refreshToken, isAuthenticated: true });
      },

      setLoading: (loading: boolean) => {
        set({ isLoading: loading });
      },

      setHasHydrated: (state: boolean) => {
        set({ hasHydrated: state });
      },
    }),
    {
      name: 'auth-storage',
      skipHydration: true,
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        state?.setHasHydrated(true);
      },
    }
  )
);
