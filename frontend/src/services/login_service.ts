import { api_client, type IApiResponse } from '@/repositories/api_client.ts';

export type ILoginResponse = IApiResponse<string>
export class LoginService {
  static async loginWithGoogle(): Promise<ILoginResponse> {
    const response = await api_client.get('/api/v1/auth/google/signin-oidc')
    return response.data;
  }
}
