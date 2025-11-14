import type IGetUserResponse from '@/infrastructure/user/IGetUserResponse.ts';
import api_client from '@/infrastructure/utils/api_client.ts';
import type { IApiResponse } from '@/infrastructure/utils/IApiResponse.ts';
import type { AxiosResponse } from 'axios';

const URLS = {
  getUserData: "/user"
}

export const UserRepository = {
  getUserData: async (): Promise<IApiResponse<IGetUserResponse>> => {
    return await api_client.get<AxiosResponse<IGetUserResponse>>(URLS.getUserData)
      .then(response => ({
        isSuccess: true,
        data: response.data.data
      }))
      .catch(error => ({
        isSuccess: false,
        ...error.response.data
      }));
  }
}
