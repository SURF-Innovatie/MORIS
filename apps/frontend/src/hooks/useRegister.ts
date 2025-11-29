import { useMutation } from '@tanstack/react-query';
import { postRegister } from '../api/generated-orval/moris';
import { RegisterRequest } from '@/api/generated-orval/model';

export const useRegister = () => {
  return useMutation({
    mutationFn: (credentials: RegisterRequest) => postRegister(credentials),
  });
};
