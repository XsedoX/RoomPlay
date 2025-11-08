import { useNotificationStore } from '@/stores/notification_store.ts';
import type ICreateRoomRequest from '@/infrastructure/models/ICreateRoomRequest.ts';
import { RoomRepository } from '@/infrastructure/repositories/room_repository.ts';
import { TSnackbarColor } from '@/infrastructure/models/TSnackbarColor.ts';
import ValidationError from '@/errors/validation_error.ts';

export const RoomService = {
  createRoom: async (roomData: ICreateRoomRequest) => {
    const response = await RoomRepository.createRoom(roomData);
    if (!response.isSuccess){
      const notificationStore = useNotificationStore();

      notificationStore.showSnackbar(response.title, TSnackbarColor.ERROR);
      if (response.validationErrors){
        throw new ValidationError(response.title, response.validationErrors);
      }
      else {
        throw new Error(response.title);
      }
    }
    return response;
  }
};
