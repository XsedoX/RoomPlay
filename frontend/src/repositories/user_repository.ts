import { type ILoginResponse, LoginService } from '@/services/login_service.ts';


export class UserRepository {
  static async loginWithGoogle(): Promise<ILoginResponse> {
    const result = await LoginService.loginWithGoogle();
    console.log(result);
    return result;
  }
}
