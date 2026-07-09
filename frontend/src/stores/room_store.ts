import { defineStore } from 'pinia';
import { useNotificationStore } from '@/stores/notification_store.ts';
import { computed, ref } from 'vue';
import type IRoomStoreModel from '@/infrastructure/room/IRoomStoreModel.ts';
import type IJoinRoomPasswordRequest from '@/infrastructure/room/IJoinRoomPasswordRequest.ts';
import type ICreateRoomRequest from '@/infrastructure/room/ICreateRoomRequest.ts';
import { TUserRole } from '@/infrastructure/user/TUserRole.ts';
import type { ISongListRoomStoreModel } from '@/infrastructure/room/ISongListRoomStoreModel.ts';
import type { TSongState } from '@/infrastructure/room/TSongState.ts';
import { Guid, type IGuid } from '@/shared/Guid.ts';
import { useRouter } from 'vue-router';
import type IPlayingSongModel from '@/infrastructure/room/IPlayingSongModel.ts';
import { TVoteStatus } from '@/infrastructure/room/TVoteStatus.ts';
import 'pinia-plugin-persistedstate';
import { RoomRepository } from '@/infrastructure/room/room_repository';
import { HttpCodes } from '@/infrastructure/utils/status_codes';
import { TSnackbarColor } from '@/infrastructure/utils/TSnackbarColor';

export const useRoomStore = defineStore(
  'room',
  () => {
    const room = ref<IRoomStoreModel | null>(null);
    const songs = ref<ISongListRoomStoreModel[] | null>(null);
    const playingSong = ref<IPlayingSongModel | null>(null);
    const router = useRouter();

    function resetRoomStore() {
      room.value = null;
      songs.value = null;
      playingSong.value = null;
    }

    async function createRoom(roomData: ICreateRoomRequest) {
      const response = await RoomRepository.createRoom(roomData);
      if (response.isSuccess) {
        await router.replace({ name: 'RoomPage' });
        return null;
      }
      if (response.validationErrors) {
        return response.validationErrors;
      }
      return null;
    }

    const isBoostAvailable = computed(() => {
      if (!room.value?.boostData) return false;
      const now = new Date();
      const elapsedSeconds = (now.getTime() - room.value.boostData.boostUsedAtUtc.getTime()) / 1000;
      return elapsedSeconds > room.value.boostData.boostCooldownSeconds;
    });
    function upVoteSong(songId: IGuid) {
      const song = songs.value?.find(
        (song) =>
          song.id.toString() === songId.toString() && song.voteStatus === TVoteStatus.notVoted,
      );
      if (song) {
        song.votes += 1;
        song.voteStatus = TVoteStatus.upvoted;
      }
    }
    function downVoteSong(songId: IGuid) {
      const song = songs.value?.find(
        (song) =>
          song.id.toString() === songId.toString() && song.voteStatus === TVoteStatus.notVoted,
      );
      if (song) {
        song.votes -= 1;
        song.voteStatus = TVoteStatus.downvoted;
      }
    }
    async function getRoom() {
      let isError = false;

      const response = await RoomRepository.getRoom();
      const notificationStore = useNotificationStore();
      if (response.isSuccess) {
        room.value = {
          name: response.data.name,
          qrCode: response.data.qrCode,
          userRole: response.data.userRole as TUserRole,
          boostData: response.data.boostData,
        };
        songs.value = response.data.songs.map((song) => ({
          ...song,
          state: song.state as TSongState,
          voteStatus: song.voteStatus as TVoteStatus,
          id: new Guid(song.id),
        }));
        playingSong.value = response.data.playingSong;
        return isError;
      }
      if (response.status === HttpCodes.notFound) {
        notificationStore.showSnackbar(
          'The room you were a member of has either expired or you got kicked or banned from it.',
          TSnackbarColor.INFO,
        );
      }
      isError = true;
      resetRoomStore();
      return isError;
    }
    async function getUserRoomMembership() {
      const response = await RoomRepository.getUserRoomMembership();
      if (!response.isSuccess) {
        resetRoomStore();
        return false;
      }
      return response.data;
    }
    async function leaveRoom() {
      await RoomRepository.leaveRoom();
      resetRoomStore();
    }
    async function joinRoomPassword(roomData: IJoinRoomPasswordRequest) {
      const response = await RoomRepository.joinRoomPassword(roomData);
      if (response.isSuccess) {
        await router.replace({ name: 'RoomPage' });
        return null;
      }
      if (response.validationErrors) {
        return response.validationErrors;
      }
      return null;
    }

    const isHost = computed(() => room.value?.userRole === TUserRole.host);

    return {
      room,
      createRoom,
      isHost,
      getRoom,
      songs,
      playingSong,
      upVoteSong,
      downVoteSong,
      getUserRoomMembership,
      leaveRoom,
      isBoostAvailable,
      joinRoomPassword,
    };
  },
  {
    persist: {
      key: 'roomplay-room',
    },
  },
);
