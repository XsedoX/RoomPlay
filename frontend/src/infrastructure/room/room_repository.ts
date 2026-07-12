import api_client from '@/infrastructure/utils/api_client.ts';
import type ICreateRoomRequest from '@/infrastructure/room/ICreateRoomRequest.ts';
import type { IGetRoomResponse } from '@/infrastructure/room/IGetRoomResponse.ts';
import type IJoinRoomPasswordRequest from './IJoinRoomPasswordRequest';
import type { IRepositoryResponse } from '../utils/IRepositoryResponse';
import { HttpCodes } from '../utils/status_codes';
import { useNotificationStore } from '@/stores/notification_store';
import { TSnackbarColor } from '../utils/TSnackbarColor';
import type { IApiProblemDetailsResponse, IApiSuccessResponse } from '../utils/IApiResponse';

const URLS = {
  createRoom: '/room',
  getRoom: '/room',
  leaveRoom: '/room',
  getUserRoomMembership: '/room/membership',
  joinRoom: '/room/join/password',
};

export const RoomRepository = {
  createRoom: async (roomData: ICreateRoomRequest): Promise<IRepositoryResponse> => {
    return await api_client
      .post<IApiSuccessResponse>(URLS.createRoom, roomData)
      .then(() => {
        return {
          isSuccess: true,
          data: undefined,
        };
      })
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        if (problemDetails.status === HttpCodes.badRequest) {
          notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
          return {
            isSuccess: false,
          };
        }
        if (problemDetails.status === HttpCodes.unprocessableEntity) {
          return {
            isSuccess: false,
            validationErrors: problemDetails.validationErrors,
          };
        }
        notificationStore.showSnackbar(problemDetails.title, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
        };
      });
  },
  getRoom: async (): Promise<IRepositoryResponse<IGetRoomResponse>> => {
    return await api_client
      .get<IApiSuccessResponse<IGetRoomResponse>>(URLS.getRoom)
      .then((response) => {
        return {
          isSuccess: true,
          data: response.data.data,
        };
      })
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        if (problemDetails.status === HttpCodes.notFound) {
          notificationStore.showSnackbar(
            'The room you were a member of has either expired or you got kicked or banned from it.',
            TSnackbarColor.INFO,
          );
          return {
            isSuccess: false,
          };
        }
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
        };
      });
  },
  getUserRoomMembership: async (): Promise<IRepositoryResponse<boolean>> => {
    return await api_client
      .get<IApiSuccessResponse<boolean>>(URLS.getUserRoomMembership)
      .then((response) => {
        return {
          isSuccess: true,
          data: response.data.data,
        };
      })
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
        };
      });
  },
  leaveRoom: async (): Promise<IRepositoryResponse> => {
    return await api_client
      .delete<IApiSuccessResponse>(URLS.leaveRoom)
      .then(() => ({
        isSuccess: true,
        data: undefined,
      }))
      .catch((_) => ({
        isSuccess: false,
      }));
  },
  joinRoomPassword: async (roomData: IJoinRoomPasswordRequest): Promise<IRepositoryResponse> => {
    return await api_client
      .put<IApiSuccessResponse>(URLS.joinRoom, roomData)
      .then((response) => ({
        isSuccess: true,
        data: response.data.data,
      }))
      .catch((error) => {
        const notificationStore = useNotificationStore();
        const problemDetails = error.response.data as IApiProblemDetailsResponse;
        notificationStore.showSnackbar(problemDetails.description, TSnackbarColor.ERROR);
        return {
          isSuccess: false,
        };
      });
  },
};
