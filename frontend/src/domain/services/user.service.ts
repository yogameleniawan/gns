import apiClient from './api-client';
import { User } from './auth.service';

/**
 * User list request parameters
 */
export interface UserListParams {
  page?: number;
  page_size?: number;
  search?: string;
  sort_by?: string;
  sort_dir?: 'asc' | 'desc';
  is_active?: boolean;
  role_id?: number;
}

/**
 * Paginated user list response
 */
export interface UserListResponse {
  items: User[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

/**
 * Create user request
 */
export interface CreateUserRequest {
  email: string;
  password: string;
  name: string;
  role_ids?: number[];
}

/**
 * Update user request (admin)
 */
export interface UpdateUserRequest {
  name?: string;
  email?: string;
  is_active?: boolean;
  avatar_url?: string;
  role_ids?: number[];
}

/**
 * User Management Service
 */
export const userService = {
  /**
   * List users with pagination and filters
   */
  async listUsers(params: UserListParams = {}): Promise<UserListResponse> {
    const queryParams = new URLSearchParams();
    
    if (params.page) queryParams.set('page', params.page.toString());
    if (params.page_size) queryParams.set('page_size', params.page_size.toString());
    if (params.search) queryParams.set('search', params.search);
    if (params.sort_by) queryParams.set('sort_by', params.sort_by);
    if (params.sort_dir) queryParams.set('sort_dir', params.sort_dir);
    if (params.is_active !== undefined) queryParams.set('is_active', params.is_active.toString());
    if (params.role_id) queryParams.set('role_id', params.role_id.toString());

    const response = await apiClient.get<UserListResponse>(`/v1/users?${queryParams.toString()}`);
    return response.data;
  },

  /**
   * Get user by ID
   */
  async getUser(id: string): Promise<User> {
    const response = await apiClient.get<User>(`/v1/users/${id}`);
    return response.data;
  },

  /**
   * Create new user
   */
  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await apiClient.post<User>('/v1/users', data);
    return response.data;
  },

  /**
   * Update user
   */
  async updateUser(id: string, data: UpdateUserRequest): Promise<User> {
    const response = await apiClient.put<User>(`/v1/users/${id}`, data);
    return response.data;
  },

  /**
   * Delete user
   */
  async deleteUser(id: string): Promise<void> {
    await apiClient.delete(`/v1/users/${id}`);
  },

  /**
   * Toggle user active status
   */
  async toggleUserStatus(id: string): Promise<User> {
    const response = await apiClient.post<User>(`/v1/users/${id}/toggle-status`);
    return response.data;
  },

  /**
   * List deleted users with pagination
   */
  async listDeletedUsers(params: UserListParams = {}): Promise<UserListResponse> {
    const queryParams = new URLSearchParams();
    
    if (params.page) queryParams.set('page', params.page.toString());
    if (params.page_size) queryParams.set('page_size', params.page_size.toString());
    if (params.search) queryParams.set('search', params.search);
    if (params.sort_by) queryParams.set('sort_by', params.sort_by);
    if (params.sort_dir) queryParams.set('sort_dir', params.sort_dir);

    const response = await apiClient.get<UserListResponse>(`/v1/users/deleted?${queryParams.toString()}`);
    return response.data;
  },

  /**
   * Restore deleted user
   */
  async restoreUser(id: string): Promise<User> {
    const response = await apiClient.post<User>(`/v1/users/${id}/restore`);
    return response.data;
  },
};
