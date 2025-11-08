import { LoginRepository } from '@/infrastructure/repositories/login_repository.ts';
import { TSnackbarColor } from '@/infrastructure/models/TSnackbarColor.ts';
import { useNotificationStore } from '@/stores/notification_store.ts';

export const LoginService = {
  logout: async () => {
    const response = await LoginRepository.logout();
    if (!response.isSuccess) {
      const notificationStore = useNotificationStore();
      notificationStore.showSnackbar(response.title, TSnackbarColor.ERROR);
    }
  },
  refreshToken: async (): Promise<boolean> => {
    const response = await LoginRepository.refreshToken();
    if (!response.isSuccess) {
      const notificationStore = useNotificationStore();

      notificationStore.showSnackbar(response.title, TSnackbarColor.ERROR);
      return false;
    }
    return true;
  },
  loginWithGoogle: async (): Promise<string|undefined> => {
    const response = await LoginRepository.loginWithGoogle();
    if (!response.isSuccess) {
      const notificationStore = useNotificationStore();

      notificationStore.showSnackbar(response.title, TSnackbarColor.ERROR);
      return undefined;
    }
    return response.data;
  }
}
