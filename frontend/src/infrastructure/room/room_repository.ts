import api_client from '@/infrastructure/utils/api_client.ts';
import type IJoinRoomPasswordRequest from './requests/IJoinRoomPasswordRequest';
import type { IRepositoryResponse } from '../utils/IRepositoryResponse';
import { HttpCodes } from '../utils/status_codes';
import { useNotificationStore } from '@/stores/notification_store';
import { TSnackbarColor } from '../utils/TSnackbarColor';
import type {
  IApiJsonPatchResponse,
  IApiProblemDetailsResponse,
  IApiSuccessResponse,
} from '../utils/IApiResponse';
import type ICreateRoomRequest from './requests/ICreateRoomRequest';
import type { IGetRoomResponse } from './responses/IGetRoomResponse';
import { useAppWebSocket } from '../websockets/useAppWebSocket';
import type { IEnqueueSongRequest } from './requests/IEnqueueSongRequest';
import { TOutgoingActionName } from '../websockets/TOutgoingActionName';
import type { TSongState } from './TSongState';
import type { TVoteStatus } from './TVoteStatus';
import { useRoomStore } from '@/stores/room_store';
import { onMounted, onUnmounted } from 'vue';
import { apiWebSocketDispatcher } from '../websockets/ApiWebSocketActionDispatcher';
import { TIncomingActionName } from '../websockets/TIncomingActionName';
import type { IPatchEnqueuedSong } from './responses/IPatchEnqueuedSong';
import type { ISongListRoomStoreModel } from './ISongListRoomStoreModel';
import { ApiError } from '../utils/ApiError';

const URLS = {
  createRoom: '/room',
  getRoom: '/room',
  leaveRoom: '/room',
  getUserRoomMembership: '/room/membership',
  joinRoom: '/room/join/password',
  upgradeToWebSocket: '/room/ws/upgrade',
  webSocketTicket: '/room/ws/ticket',
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
      .catch((error: ApiError) => {
        if (error.problemDetails === undefined) {
          return {
            isSuccess: false,
          };
        }
        switch (error.problemDetails.status) {
          case HttpCodes.unprocessableEntity:
            return { isSuccess: false, validationErrors: error.problemDetails.validationErrors };
          default:
            return { isSuccess: false };
        }
      });
  },
  getRoom: async (): Promise<IRepositoryResponse<IGetRoomResponse>> => {
    return await api_client
      .get<IApiSuccessResponse<IGetRoomResponse>>(URLS.getRoom, { silent: true })
      .then((response) => {
        return {
          isSuccess: true,
          data: response.data.data,
        };
      })
      .catch((error: ApiError) => {
        const notificationStore = useNotificationStore();
        if (error.problemDetails === undefined) {
          notificationStore.showSnackbar(
            'Sorry, something went wrong. Please try again later.',
            TSnackbarColor.ERROR,
          );
          return {
            isSuccess: false,
          };
        }
        switch (error.problemDetails.status) {
          case HttpCodes.notFound:
            notificationStore.showSnackbar(
              'The room you were a member of has either expired or you got kicked or banned from it.',
              TSnackbarColor.INFO,
            );
            return {
              isSuccess: false,
            };
          default:
            notificationStore.showSnackbar(error.problemDetails.description, TSnackbarColor.ERROR);
            return {
              isSuccess: false,
            };
        }
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
      .catch(() => {
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
      .catch(() => {
        return {
          isSuccess: false,
        };
      });
  },
  useRoomWebSocket: () => {
    const { send, close, reconnectAttempt, reconnectDelay, status } = useAppWebSocket(
      URLS.upgradeToWebSocket,
      URLS.webSocketTicket,
    );
    const notificationStore = useNotificationStore();
    const roomStore = useRoomStore();

    const enqueueSong = (request: IEnqueueSongRequest) =>
      send(TOutgoingActionName.enqueueSong, request);

    const handleSongEnqueuedError = (response: IApiProblemDetailsResponse) => {
      notificationStore.showSnackbar(response.title, TSnackbarColor.ERROR);
    };

    const handleEnqueuedSongPatch = (response: IApiJsonPatchResponse<IPatchEnqueuedSong>) => {
      const patch = response.data;
      patch.forEach((operation) => {
        const song: ISongListRoomStoreModel = {
          addedBy: operation.value.addedBy!,
          votes: operation.value.votes!,
          albumCoverUrl: operation.value.albumCoverUrl!,
          title: operation.value.title!,
          author: operation.value.author!,
          id: operation.value.id!,
          state: operation.value.state as TSongState,
          voteStatus: operation.value.voteStatus as TVoteStatus,
        };
        roomStore.enqueueSong(song);
      });
    };
    let unsubFunctions: Array<() => void> = [];
    onMounted(() => {
      unsubFunctions = [
        apiWebSocketDispatcher.on(TIncomingActionName.songEnqueuedError, handleSongEnqueuedError),
        apiWebSocketDispatcher.on(TIncomingActionName.enqueuedSongPatch, handleEnqueuedSongPatch),
      ];
    });

    onUnmounted(() => {
      unsubFunctions.forEach((unsub) => unsub());
      unsubFunctions = [];
    });
    return { enqueueSong, close, reconnectAttempt, reconnectDelay, status };
  },
};
