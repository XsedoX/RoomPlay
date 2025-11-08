import api_client from '@/infrastructure/repositories/api_client.ts';
import type ICreateRoomRequest from '@/infrastructure/models/ICreateRoomRequest.ts';
import type { IApiResponse } from '@/infrastructure/models/IApiResponse.ts';
import type { IGuid } from '@/shared/Guid.ts';

const URLS = {
  createRoom: "/room"
}

export const RoomRepository = {
  createRoom: async (roomData: ICreateRoomRequest): Promise<IApiResponse<IGuid>> => {
   return api_client.post(URLS.createRoom, roomData)
     .then(response => ({
       isSuccess: true,
       data: response.data
     }))
     .catch(error => ({
       isSuccess: false,
       ...error.response.data
     }));
  }
}
