import { defineStore } from 'pinia';
import { ref } from 'vue';
import type IRoomStoreModel from '@/infrastructure/models/IRoomStoreModel.ts';
import type ICreateRoomRequest from '@/infrastructure/models/ICreateRoomRequest.ts';
import { RoomService } from '@/infrastructure/services/room_service.ts';
import ValidationError from '@/errors/validation_error.ts'

export const useRoomStore = defineStore('room', () => {
  const room = ref<IRoomStoreModel|null>()

  async function createRoom(roomData: ICreateRoomRequest) {
    return await RoomService.createRoom(roomData)
      .then(r => {
        room.value = {
          id: r.data,
          name: roomData.name
        }
        return null
      })
      .catch(err => {
        if (err instanceof ValidationError) {
          return err.fieldErrors
        }
        return null
      });
  }

  return { room, createRoom };
},
  {
    persist: {
      key: 'roomplay-room',
      }
  }
);
