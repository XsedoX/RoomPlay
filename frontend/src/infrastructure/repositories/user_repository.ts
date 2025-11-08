import type IUserDataResponse from '@/infrastructure/models/IUserDataResponse.ts';
import api_client from '@/infrastructure/repositories/api_client.ts';
import type { IApiResponse } from '@/infrastructure/models/IApiResponse.ts';
import type { AxiosResponse } from 'axios';

const URLS = {
  getUserData: "/user"
}

export const UserRepository = {
  getUserData: async (): Promise<IApiResponse<IUserDataResponse>> => {
    return await api_client.get<AxiosResponse<IUserDataResponse>>(URLS.getUserData)
      .then(response => ({
        isSuccess: true,
        data: response.data.data
      }))
      .catch(error => ({
        isSuccess: false,
        ...error.response.data.data
      }));
  }
}
