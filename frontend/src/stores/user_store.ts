import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type { ILoggedInUserStoreModel } from '@/infrastructure/user/ILoggedInUserStoreModel.ts';
import { UserRepository } from '@/infrastructure/user/user_repository';
import { useNotificationStore } from './notification_store';
import { TSnackbarColor } from '@/infrastructure/utils/TSnackbarColor';
import { AuthenticationRepository } from '@/infrastructure/authentication/authentication_repository';

export const useUserStore = defineStore(
  'user',
  () => {
    const user = ref<ILoggedInUserStoreModel | null>(null);

    async function getUserData() {
      const userDataResponse = await UserRepository.getUserData();
      if (userDataResponse.isSuccess) {
        user.value = {
          name: userDataResponse.data.name,
          surname: userDataResponse.data.surname,
        };
        return null;
      }
      const notificationStore = useNotificationStore();
      notificationStore.showSnackbar(userDataResponse.title, TSnackbarColor.ERROR);
      return null;
    }

    async function logout() {
      await AuthenticationRepository.logout();
      user.value = null;
    }
    const usersInitials = computed(() => {
      return `${user.value?.name[0]}${user.value?.surname[0]}`.toUpperCase();
    });

    return { user, getUserData, logout, usersInitials };
  },
  {
    persist: {
      key: 'roomplay-user',
    },
  },
);
