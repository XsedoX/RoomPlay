import { AuthenticationRepository } from '@/infrastructure/authentication/authentication_repository.ts';
import { TSnackbarColor } from '@/infrastructure/utils/TSnackbarColor.ts';
import { useNotificationStore } from '@/stores/notification_store.ts';

export const AuthenticationService = {
  logout: async () => {
    await AuthenticationRepository.logout();
  },
  refreshToken: async (): Promise<boolean> => {
    const response = await AuthenticationRepository.refreshToken();
    return response.isSuccess;
  },
  loginWithGoogle: async (): Promise<string | undefined> => {
    const response = await AuthenticationRepository.loginWithGoogle();
    if (!response.isSuccess) {
      const notificationStore = useNotificationStore();

      notificationStore.showSnackbar(response.title, TSnackbarColor.ERROR);
      return undefined;
    }
    return response.data;
  },
};
