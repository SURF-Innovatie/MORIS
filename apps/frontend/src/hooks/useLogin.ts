import { useMutation } from '@tanstack/react-query';
import { useAuth } from './useAuth';
import { postLogin } from '../api/generated-orval/moris';
import { LoginRequest } from '../api/generated-orval/model';
import { useNavigate } from 'react-router-dom';

export const useLogin = () => {
  const { login } = useAuth();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (credentials: LoginRequest) => postLogin(credentials),
    onSuccess: (data) => {
      // Store token and user info
      if (data.token && data.user) {
        login(data.token, data.user);
        // Navigate to dashboard or home
        navigate('/dashboard');
      }
    },
  });
};
