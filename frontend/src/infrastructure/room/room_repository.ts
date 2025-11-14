import api_client from '@/infrastructure/utils/api_client.ts';
import type ICreateRoomRequest from '@/infrastructure/room/ICreateRoomRequest.ts';
import type { IApiResponse } from '@/infrastructure/utils/IApiResponse.ts';
import { type AxiosResponse } from 'axios';
import type { IGetRoomResponse } from '@/infrastructure/room/IGetRoomResponse.ts';

const URLS = {
  createRoom: "/room",
  getRoom: "/room",
  leaveRoom: "/room",
  getUserRoomMembership: "/room/membership"
}

export const RoomRepository = {
  createRoom: async (roomData: ICreateRoomRequest): Promise<IApiResponse> => {
   return await api_client.post<AxiosResponse>(URLS.createRoom, roomData)
     .then(response => ({
       isSuccess: true,
       data: response.data.data
     }))
     .catch(error => ({
       isSuccess: false,
       ...error.response.data
     }));
  },
  getRoom: async (): Promise<IApiResponse<IGetRoomResponse>> => {
    return await api_client.get<AxiosResponse<IGetRoomResponse>>(URLS.getRoom)
      .then(response => ({
        isSuccess: true,
        data: response.data.data
      }))
      .catch(error => {
        console.log(error)
        return {
        isSuccess: false,
        ...error.response.data
        }
      })
  },
  getUserRoomMembership: async (): Promise<IApiResponse<boolean>> => {
    return await api_client.get<AxiosResponse<boolean>>(URLS.getUserRoomMembership)
      .then(response => ({
        isSuccess: true,
        data: response.data.data
      }))
      .catch(error => ({
        isSuccess: false,
        ...error.response.data
      }));
  },
  leaveRoom: async (): Promise<IApiResponse> => {
    return await api_client.delete<AxiosResponse>(URLS.leaveRoom)
      .then(response => ({
        isSuccess: true,
        data: response.data.data
      }))
      .catch(error => ({
        isSuccess: false,
        ...error.response.data
      }))
  }

}
