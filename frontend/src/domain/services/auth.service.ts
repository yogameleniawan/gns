import apiClient from './api-client';

/**
 * User interface matching backend response
 */
export interface User {
  id: string;
  email: string;
  name: string;
  avatar_url?: string;
  is_oauth: boolean;
  oauth_provider?: string;
  is_active: boolean;
  email_verified: boolean;
  roles?: string[];
  permissions?: string[];
  created_at: string;
}

/**
 * Authentication response from backend
 */
export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

/**
 * Register request payload
 */
export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

/**
 * Login request payload
 */
export interface LoginRequest {
  email: string;
  password: string;
}

/**
 * Google OAuth request payload
 */
export interface GoogleOAuthRequest {
  id_token: string;
}

/**
 * Refresh token request
 */
export interface RefreshTokenRequest {
  refresh_token: string;
}

/**
 * Update profile request
 */
export interface UpdateProfileRequest {
  name?: string;
  avatar_url?: string;
}

/**
 * Change password request
 */
export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

/**
 * Profile response
 */
export interface ProfileResponse {
  user: User;
}

/**
 * Authentication Service
 */
export const authService = {
  /**
   * Register a new user
   */
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/v1/auth/register', data);
    return response.data;
  },

  /**
   * Login with email and password
   */
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/v1/auth/login', data);
    return response.data;
  },

  /**
   * Login with Google OAuth
   */
  async loginWithGoogle(data: GoogleOAuthRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/v1/auth/oauth/google', data);
    return response.data;
  },

  /**
   * Refresh access token
   */
  async refreshToken(data: RefreshTokenRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/v1/auth/refresh', data);
    return response.data;
  },

  /**
   * Logout (invalidate session)
   */
  async logout(): Promise<void> {
    await apiClient.post('/v1/auth/logout');
  },

  /**
   * Get user profile
   */
  async getProfile(): Promise<ProfileResponse> {
    const response = await apiClient.get<ProfileResponse>('/v1/auth/profile');
    return response.data;
  },

  /**
   * Update user profile
   */
  async updateProfile(data: UpdateProfileRequest): Promise<User> {
    const response = await apiClient.put<User>('/v1/auth/profile', data);
    return response.data;
  },

  /**
   * Change password
   */
  async changePassword(data: ChangePasswordRequest): Promise<void> {
    await apiClient.post('/v1/auth/change-password', data);
  },
};
