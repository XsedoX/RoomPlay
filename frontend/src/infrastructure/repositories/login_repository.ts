import api_client from '@/infrastructure/repositories/api_client.ts';
import type IApiResponse from '@/infrastructure/models/IApiResponse.ts'
import type IUserDataResponse from '@/infrastructure/models/IUserDataResponse.ts'

const URLS = {
  loginWithGoogle: '/auth/google/signin-oidc',
  getUserData: "/user"
}

export const LoginRepository = {
  loginWithGoogle: async (): Promise<string> => {
    const response = await api_client.get<IApiResponse<string>>(URLS.loginWithGoogle);
    console.log(response);
    return response.data.data as string;
  },
  userData: async (): Promise<IUserDataResponse> => {
    const response = await api_client.get<IApiResponse<IUserDataResponse>>(URLS.getUserData);
    console.log(response.data.data);
    return response.data.data as IUserDataResponse;
  }
}
