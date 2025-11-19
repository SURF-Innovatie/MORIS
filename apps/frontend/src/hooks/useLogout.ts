import { useCallback } from 'react';
import { useAuth } from './useAuth';
import { useNavigate } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';

export const useLogout = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  return useCallback(() => {
    logout();
    // Clear all cached queries
    queryClient.clear();
    // Navigate to login
    navigate('/login');
  }, [logout, navigate, queryClient]);
};
