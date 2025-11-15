import type ICreateRoomRequest from '@/infrastructure/room/ICreateRoomRequest.ts';
import { RoomRepository } from '@/infrastructure/room/room_repository.ts';
import ValidationError from '@/errors/validation_error.ts';
import type { IGetRoomResponse } from '@/infrastructure/room/IGetRoomResponse.ts';
import { HttpCodes } from '@/infrastructure/utils/status_codes.ts';
import NotFoundError from '@/errors/not_found_error.ts';
import { useNotificationStore } from '@/stores/notification_store.ts';
import { TSnackbarColor } from '@/infrastructure/utils/TSnackbarColor.ts';

export const RoomService = {
  createRoom: async (roomData: ICreateRoomRequest) => {
    const response = await RoomRepository.createRoom(roomData);
    if (!response.isSuccess) {
      if (response.validationErrors) {
        throw new ValidationError(response.title, response.validationErrors);
      } else {
        throw new Error(response.title);
      }
    }
    return response;
  },
  getRoom: async (): Promise<IGetRoomResponse> => {
    const response = await RoomRepository.getRoom();
    const notificationStore = useNotificationStore();
    if (!response.isSuccess) {
      if (response.validationErrors) {
        throw new ValidationError(response.title, response.validationErrors);
      } else if (response.status === HttpCodes.notFound) {
        notificationStore.showSnackbar(
          'The room you were a member of has either expired or you get kicked or banned from it.',
          TSnackbarColor.INFO,
        );
        throw new NotFoundError(response.title);
      } else {
        notificationStore.showSnackbar(response.title, TSnackbarColor.ERROR);
        console.error('CUSTOM', response);
        throw new Error(response.title);
      }
    }
    return response.data;
  },
  getUserRoomMembership: async () => {
    const response = await RoomRepository.getUserRoomMembership();
    if (!response.isSuccess) {
      return false;
    }
    return response.data;
  },
  leaveRoom: async () => {
    const response = await RoomRepository.leaveRoom();
    if (!response.isSuccess) {
      if (response.validationErrors) {
        throw new ValidationError(response.title, response.validationErrors);
      } else {
        throw new Error(response.title);
      }
    }
    return response;
  },
};
