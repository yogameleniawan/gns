import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userService, UserListParams, CreateUserRequest, UpdateUserRequest } from '@/src/domain/services/user.service';

// Query keys
export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (params: UserListParams) => [...userKeys.lists(), params] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
};

/**
 * Hook to fetch users with pagination
 */
export function useUsers(params: UserListParams = {}) {
  return useQuery({
    queryKey: userKeys.list(params),
    queryFn: () => userService.listUsers(params),
    staleTime: 30000, // 30 seconds
  });
}

/**
 * Hook to fetch single user
 */
export function useUser(id: string) {
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: () => userService.getUser(id),
    enabled: !!id,
  });
}

/**
 * Hook to create user
 */
export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserRequest) => userService.createUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
}

/**
 * Hook to update user
 */
export function useUpdateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserRequest }) =>
      userService.updateUser(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
      queryClient.invalidateQueries({ queryKey: userKeys.detail(variables.id) });
    },
  });
}

/**
 * Hook to delete user
 */
export function useDeleteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => userService.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
}

/**
 * Hook to toggle user status
 */
export function useToggleUserStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => userService.toggleUserStatus(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
      queryClient.invalidateQueries({ queryKey: userKeys.detail(id) });
    },
  });
}

/**
 * Hook to fetch deleted users with pagination
 */
export function useDeletedUsers(params: UserListParams = {}) {
  return useQuery({
    queryKey: [...userKeys.all, 'deleted', params] as const,
    queryFn: () => userService.listDeletedUsers(params),
    staleTime: 30000,
  });
}

/**
 * Hook to restore user
 */
export function useRestoreUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => userService.restoreUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
      queryClient.invalidateQueries({ queryKey: [...userKeys.all, 'deleted'] });
    },
  });
}
