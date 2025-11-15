import type { ILoggedInUserStoreModel } from '@/infrastructure/user/ILoggedInUserStoreModel.ts';
import { UserRepository } from '@/infrastructure/user/user_repository.ts';
import { useNotificationStore } from '@/stores/notification_store.ts';
import { TSnackbarColor } from '@/infrastructure/utils/TSnackbarColor.ts';
import ValidationError from '@/errors/validation_error.ts';

export const UserService = {
  getUserData: async (): Promise<ILoggedInUserStoreModel> => {
    const userDataResponse = await UserRepository.getUserData();
    if (!userDataResponse.isSuccess) {
      const notificationStore = useNotificationStore();

      notificationStore.showSnackbar(userDataResponse.title, TSnackbarColor.ERROR);
      if (userDataResponse.validationErrors) {
        throw new ValidationError(userDataResponse.title, userDataResponse.validationErrors);
      } else {
        throw new Error(userDataResponse.title);
      }
    }
    return {
      name: userDataResponse.data.name,
      surname: userDataResponse.data.surname,
    };
  },
};
