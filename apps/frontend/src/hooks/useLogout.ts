import { useMutation } from '@tanstack/react-query';
import { useAuth } from './useAuth';
import { useQueryClient } from '@tanstack/react-query';

export const useLogout = () => {
  const { logout } = useAuth();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      // No API call needed, just clear local state
      return Promise.resolve();
    },
    onSuccess: () => {
      logout();
      // Clear all cached queries
      queryClient.clear();
    },
  });
};
