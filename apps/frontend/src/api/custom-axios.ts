import Axios, { AxiosRequestConfig, AxiosError } from "axios";

export const AXIOS_INSTANCE = Axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "/api",
});

export interface BackendError {
  code: string;
  message: string;
  details?: any;
  status: number;
  timestamp?: string;
}

import { STORAGE_KEY_AUTH_TOKEN } from "@/lib/constants";

// Request interceptor to add JWT token
AXIOS_INSTANCE.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem(STORAGE_KEY_AUTH_TOKEN);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
AXIOS_INSTANCE.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response) {
      // Check if 401 and dispatch session expired event
      if (error.response.status === 401) {
        // Do not clear token immediately, let the user decide when to log out via the dialog
        window.dispatchEvent(new CustomEvent("auth:session-expired"));
      }

      const backendError = error.response.data as Partial<BackendError>;
      const enrichedError: BackendError = {
        code: backendError.code || "UNKNOWN_ERROR",
        message: backendError.message || "An unexpected error occurred.",
        details: backendError.details,
        status: error.response.status,
        timestamp: backendError.timestamp,
      };
      return Promise.reject(enrichedError);
    }
    const networkError: BackendError = {
      code: "NETWORK_ERROR",
      message: error.message || "Network error or server unreachable.",
      status: error.status || 0,
    };
    return Promise.reject(networkError);
  }
);

export const customInstance = async <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig
): Promise<T> => {
  const source = Axios.CancelToken.source();
  const promise = AXIOS_INSTANCE({
    ...config,
    ...options,
    cancelToken: source.token,
  }).then((res) => res.data);
  // @ts-ignore
  promise.cancel = () => {
    source.cancel("Query was cancelled");
  };
  return promise;
};

export default customInstance;
