import api_client from '@/infrastructure/utils/api_client.ts';
import type { IRepositoryResponse } from '../utils/IRepositoryResponse';
import type { IApiProblemDetailsResponse, IApiSuccessResponse } from '../utils/IApiResponse';
import { useNotificationStore } from '@/stores/notification_store';
import { TSnackbarColor } from '../utils/TSnackbarColor';

const URLS = {
  loginWithGoogle: '/auth/google/signin-oidc',
  logout: '/auth/logout',
  refreshToken: '/auth/refresh-token',
};

export const AuthenticationRepository = {
  loginWithGoogle: async (): Promise<IRepositoryResponse<string>> => {
    return await api_client
      .get<IApiSuccessResponse<string>>(URLS.loginWithGoogle)
      .then((response) => ({
        isSuccess: true,
        data: response.data.data,
      }))
      .catch((error) => {
        // Assuming error.response.data contains the failure details
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
        };
      });
  },

  logout: async (): Promise<IRepositoryResponse> => {
    return await api_client
      .post<IApiSuccessResponse>(URLS.logout)
      .then(() => ({
        isSuccess: true,
        data: undefined,
      }))
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
          ...error.response.data,
        };
      });
  },

  refreshToken: async (): Promise<IRepositoryResponse> => {
    return await api_client
      .post<IApiSuccessResponse>(URLS.refreshToken)
      .then(() => ({
        isSuccess: true,
        data: undefined,
      }))
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
          ...error.response.data,
        };
      });
  },
};
