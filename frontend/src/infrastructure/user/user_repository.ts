import type IGetUserResponse from '@/infrastructure/user/IGetUserResponse.ts';
import api_client from '@/infrastructure/utils/api_client.ts';
import type { IRepositoryResponse } from '../utils/IRepositoryResponse';
import type { IApiSuccessResponse } from '../utils/IApiResponse';

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
      .catch(() => {
        return {
          isSuccess: false,
        };
      });
  },
};
