import type { ILoggedInUserStoreModel } from '@/infrastructure/models/ILoggedInUserStoreModel.ts';
import { Guid } from '@/shared/Guid.ts';
import type { TUserRole } from '@/infrastructure/models/TUserRole.ts';
import { UserRepository } from '@/infrastructure/repositories/user_repository.ts';
import { useNotificationStore } from '@/stores/notification_store.ts';
import { TSnackbarColor } from '@/infrastructure/models/TSnackbarColor.ts';
import ValidationError from '@/errors/validation_error.ts';

export const UserService = {
  getUserData: async (): Promise<ILoggedInUserStoreModel> => {
    const userDataResponse = await UserRepository.getUserData();
    if (!userDataResponse.isSuccess) {
      const notificationStore = useNotificationStore();

      notificationStore.showSnackbar(userDataResponse.title, TSnackbarColor.ERROR);
      if (userDataResponse.validationErrors){
        throw new ValidationError(userDataResponse.title, userDataResponse.validationErrors);
      }
      else {
        throw new Error(userDataResponse.title);
      }
    }
    return {
      name: userDataResponse.data.name,
      surname: userDataResponse.data.surname,
      roomId: userDataResponse.data.roomId === null ? null : new Guid(userDataResponse.data.roomId),
      role: userDataResponse.data.role === null ? null: userDataResponse.data.role as TUserRole
    }
  }
}
