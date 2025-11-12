import Axios, { AxiosRequestConfig, AxiosError } from 'axios';

export const AXIOS_INSTANCE = Axios.create({ baseURL: import.meta.env.VITE_API_BASE_URL || '/api' });

export interface BackendError {
  code: string;
  message: string;
  details?: any;
  status: number;
  timestamp?: string;
}

export const customInstance = async <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig,
): Promise<T> => {
  const source = Axios.CancelToken.source();
  const promise = AXIOS_INSTANCE({ ...config, ...options, cancelToken: source.token }).then(
    (res) => res.data,
  );
  // @ts-ignore
  promise.cancel = () => {
    source.cancel('Query was cancelled');
  };
  return promise;
};

AXIOS_INSTANCE.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response) {
      const backendError = error.response.data as Partial<BackendError>;
      const enrichedError: BackendError = {
        code: backendError.code || 'UNKNOWN_ERROR',
        message: backendError.message || 'An unexpected error occurred.',
        details: backendError.details,
        status: error.response.status,
        timestamp: backendError.timestamp,
      };
      return Promise.reject(enrichedError);
    }
    const networkError: BackendError = {
      code: 'NETWORK_ERROR',
      message: error.message || 'Network error or server unreachable.',
      status: error.status || 0,
    };
    return Promise.reject(networkError);
  },
);

export default customInstance;