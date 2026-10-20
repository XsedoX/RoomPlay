import axios, { type AxiosInstance } from 'axios';
import { useUserStore } from '@/stores/user_store.ts';
import { PlatformDiscoverer } from '@/infrastructure/utils/platform_discoverer.ts';
import router from '@/router';
import { AuthenticationRepository } from '../authentication/authentication_repository';
import type { IApiProblemDetailsResponse } from './IApiResponse';
import { ApiError } from './ApiError';
import { useNotificationStore } from '@/stores/notification_store';
import { TSnackbarColor } from './TSnackbarColor';

declare module 'axios' {
  export interface AxiosRequestConfig {
    //NOTE: Skip the global error snackbar when true. Defaults to false.
    silent?: boolean;
  }
}
const axiosParams = {
  baseURL: import.meta.env['VITE_API_BASE_URL'],
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
    'X-Device-Type': PlatformDiscoverer.getDeviceType(),
  },
  timeout: 10000,
  withCredentials: true,
};

function normalize(error: unknown): ApiError {
  if (axios.isAxiosError(error)) {
    const status = error.response?.status;
    const problemDetails = error.response?.data as IApiProblemDetailsResponse | undefined;
    return new ApiError(
      error.message,
      status,
      status !== undefined && error.response?.data ? problemDetails : undefined,
      error,
    );
  }
  // non-axios: includes the TypeError from inside a catch and any thrown boilerplate
  return new ApiError('Unknown API error', undefined, undefined, error);
}
const api_client = axios.create(axiosParams);

const api = (axiosInstance: AxiosInstance) => {
  return {
    get: async <T>(url: string, config = {}) => {
      try {
        return await axiosInstance.get<T>(url, config);
      } catch (error) {
        throw normalize(error);
      }
    },
    post: async <T>(url: string, body?: unknown, config = {}) => {
      try {
        return await axiosInstance.post<T>(url, body, config);
      } catch (error) {
        throw normalize(error);
      }
    },
    delete: async <T>(url: string, config = {}) => {
      try {
        return await axiosInstance.delete<T>(url, config);
      } catch (error) {
        throw normalize(error);
      }
    },
    put: async <T>(url: string, config = {}) => {
      try {
        return await axiosInstance.put<T>(url, config);
      } catch (error) {
        throw normalize(error);
      }
    },
  };
};
async function forceLogout() {
  const userStore = useUserStore();
  await userStore.logout();
  await router.replace({ name: 'LoginPage' });
}

function createRefreshTokenInterceptor() {
  const interceptor = api_client.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config;
      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error);
      }
      if (error.response?.status === 401 && error.response?.data) {
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        if (
          problemDetails.type ===
          'ExternalAuthenticationService.RefreshAccessTokenWithExternalProvider.RefreshTokenExpired'
        ) {
          await forceLogout();
          return Promise.reject(error);
        }
      }
      originalRequest._retry = true;
      api_client.interceptors.response.eject(interceptor);

      try {
        const response = await AuthenticationRepository.refreshToken();
        if (!response.isSuccess) {
          await forceLogout();
          return Promise.reject(error);
        }
        return api_client(originalRequest);
      } catch (error) {
        return Promise.reject(error);
      } finally {
        createRefreshTokenInterceptor();
      }
    },
  );
}

createRefreshTokenInterceptor();

api_client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const config = error.config ?? {};
    if (!config.silent && error.response?.status !== 401) {
      const pd = error.response?.data as IApiProblemDetailsResponse | undefined;
      useNotificationStore().showSnackbar(
        pd?.title ?? 'Unexpected error. Please try again.',
        TSnackbarColor.ERROR,
      );
    }
    return Promise.reject(error); // still raw — the wrapper normalizes once, below
  },
);
export default api(api_client);
