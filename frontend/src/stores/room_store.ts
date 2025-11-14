import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type IRoomStoreModel from '@/infrastructure/room/IRoomStoreModel.ts';
import { RoomService } from '@/infrastructure/room/room_service.ts';
import type ICreateRoomRequest from '@/infrastructure/room/ICreateRoomRequest.ts';
import ValidationError from '@/errors/validation_error.ts';
import { TUserRole } from '@/infrastructure/user/TUserRole.ts';
import type { ISongListRoomStoreModel } from '@/infrastructure/room/ISongListRoomStoreModel.ts';
import type { TSongState } from '@/infrastructure/room/TSongState.ts';
import { Guid, type IGuid } from '@/shared/Guid.ts';
import { useRouter } from 'vue-router';
import type IPlayingSongModel from '@/infrastructure/room/IPlayingSongModel.ts';
import type { TVoteStatus } from '@/infrastructure/room/TVoteStatus.ts';
import 'pinia-plugin-persistedstate';

export const useRoomStore = defineStore('room', () => {
  const room = ref<IRoomStoreModel|null>(null)
  const songs = ref<ISongListRoomStoreModel[]|null>(null)
  const playingSong = ref<IPlayingSongModel|null>(null)
  const router = useRouter();

  function resetRoomStore() {
    room.value = null;
    songs.value = null;
    playingSong.value = null;
  }

  async function createRoom(roomData: ICreateRoomRequest) {
    return await RoomService.createRoom(roomData)
      .then(async _ => {
        await router.replace({ name: 'RoomPage'})
        return null
      })
      .catch(err => {
        if (err instanceof ValidationError) {
          return err.fieldErrors
        }
        return null
      });
  }
  const isBoostAvailable = computed(()=> {
    if (!room.value?.boostData) return false
    const now = new Date();
    const elapsedSeconds = (now.getTime() - room.value.boostData.boostUsedAtUtc.getTime()) / 1000;
    return elapsedSeconds > room.value.boostData.boostCooldownSeconds
  })
  function upVoteSong(songId: IGuid) {
    const song = songs.value?.find((song) =>
      song.id.toString() === songId.toString());
    if (song) {
      song.votes += 1;
    }
  }
  function downVoteSong(songId: IGuid) {
    const song = songs.value?.find((song) => song.id.toString() === songId.toString());
    if (song && song.votes > 0) {
      song.votes -= 1;
    }
  }
  async function getRoom() {
    let isError = false
    await RoomService.getRoom()
      .then(response => {
        room.value = {
          name: response.name,
          qrCode: response.qrCode,
          userRole: response.userRole as TUserRole,
          boostData: response.boostData
        };
        songs.value = response.songs.map(song => ({
          ...song,
          state: song.state as TSongState,
          voteStatus: song.voteStatus as TVoteStatus,
          id: new Guid(song.id)
        }))
        playingSong.value = response.playingSong
      })
      .catch(async _ => {
        isError = true
        resetRoomStore()
      })
    return isError
  }
  async function getUserRoomMembership() {
    return await RoomService.getUserRoomMembership()
      .then(isUserInRoom => {
        if (!isUserInRoom) {
          resetRoomStore()
          return false
        }
        return true
      })
  }
  async function leaveRoom() {
    await RoomService.leaveRoom()
      .then(() => {})
      .catch(() => {})
      .finally(async () => {
        resetRoomStore()
      });
  }
  const isHost = computed(() => room.value?.userRole === TUserRole.host);

  return { room,
    createRoom,
    isHost,
    getRoom,
    songs,
    playingSong,
    upVoteSong,
    downVoteSong,
    getUserRoomMembership,
    leaveRoom,
    isBoostAvailable
  };
},
  {
    persist: {
      key: 'roomplay-room',
      }
  }
);
