import api_client from '@/infrastructure/repositories/api_client.ts';
import type { IApiResponse } from '@/infrastructure/models/IApiResponse.ts';
import type { AxiosResponse } from 'axios';

const URLS = {
  loginWithGoogle: '/auth/google/signin-oidc',
  logout: '/auth/logout',
  refreshToken: '/auth/refresh-token'
}

export const LoginRepository = {
  loginWithGoogle: async (): Promise<IApiResponse<string>> => {
    return api_client.get<AxiosResponse<string>>(URLS.loginWithGoogle)
      .then(response => ({
        isSuccess: true,
        data: response.data.data
      }))
      .catch(error => {
        // Assuming error.response.data contains the failure details
        return {
          isSuccess: false,
          ...error.response.data
        };
      });
  },

  logout: async (): Promise<IApiResponse> => {
    return api_client.post<AxiosResponse<void>>(URLS.logout)
      .then(() => ({
        isSuccess: true,
        data: undefined
      }))
      .catch(error => ({
        isSuccess: false,
        ...error.response.data.data
      }));
  },

  refreshToken: async (): Promise<IApiResponse> => {
    return api_client.post<AxiosResponse<void>>(URLS.refreshToken)
      .then(() => ({
        isSuccess: true,
        data: undefined
      }))
      .catch(error => ({
        isSuccess: false,
        ...error.response.data.data
      }));
  }
};
