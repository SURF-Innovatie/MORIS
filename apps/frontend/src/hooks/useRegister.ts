import { useMutation } from "@tanstack/react-query";
import { postRegister } from "@api/moris";
import { RegisterRequest } from "@api/model";

export const useRegister = () => {
  return useMutation({
    mutationFn: (credentials: RegisterRequest) => postRegister(credentials),
  });
};
