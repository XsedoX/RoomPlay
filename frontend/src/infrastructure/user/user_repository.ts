import type IGetUserResponse from '@/infrastructure/user/IGetUserResponse.ts';
import api_client from '@/infrastructure/utils/api_client.ts';
import type { IRepositoryResponse } from '../utils/IRepositoryResponse';
import type { IApiProblemDetailsResponse, IApiSuccessResponse } from '../utils/IApiResponse';
import { useNotificationStore } from '@/stores/notification_store';
import { TSnackbarColor } from '../utils/TSnackbarColor';

const URLS = {
  getUserData: '/user',
};

export const UserRepository = {
  getUserData: async (): Promise<IRepositoryResponse<IGetUserResponse>> => {
    return await api_client
      .get<IApiSuccessResponse<IGetUserResponse>>(URLS.getUserData)
      .then((response) => ({
        isSuccess: true,
        data: response.data.data,
      }))
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response?.data as IApiProblemDetailsResponse;
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
        };
      });
  },
};
