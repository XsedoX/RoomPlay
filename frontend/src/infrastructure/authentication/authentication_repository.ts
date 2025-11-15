import api_client from '@/infrastructure/utils/api_client.ts';
import type { IApiResponse } from '@/infrastructure/utils/IApiResponse.ts';
import type { AxiosResponse } from 'axios';

const URLS = {
  loginWithGoogle: '/auth/google/signin-oidc',
  logout: '/auth/logout',
  refreshToken: '/auth/refresh-token',
};

export const AuthenticationRepository = {
  loginWithGoogle: async (): Promise<IApiResponse<string>> => {
    return await api_client
      .get<AxiosResponse<string>>(URLS.loginWithGoogle)
      .then((response) => ({
        isSuccess: true,
        data: response.data.data,
      }))
      .catch((error) => {
        // Assuming error.response.data contains the failure details
        return {
          isSuccess: false,
          ...error.response.data,
        };
      });
  },

  logout: async (): Promise<IApiResponse> => {
    return await api_client
      .post<AxiosResponse<void>>(URLS.logout)
      .then(() => ({
        isSuccess: true,
        data: undefined,
      }))
      .catch((error) => ({
        isSuccess: false,
        ...error.response.data,
      }));
  },

  refreshToken: async (): Promise<IApiResponse> => {
    return await api_client
      .post<AxiosResponse<void>>(URLS.refreshToken)
      .then(() => ({
        isSuccess: true,
        data: undefined,
      }))
      .catch((error) => ({
        isSuccess: false,
        ...error.response.data,
      }));
  },
};
