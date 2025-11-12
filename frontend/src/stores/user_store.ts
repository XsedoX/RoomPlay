import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type { ILoggedInUserStoreModel } from '@/infrastructure/user/ILoggedInUserStoreModel.ts';
import { UserService } from '@/infrastructure/user/user_service.ts';
import { AuthenticationService } from '@/infrastructure/authentication/authentication_service.ts';

export const useUserStore = defineStore('user',  () => {
  const user = ref<ILoggedInUserStoreModel | null>(null);

  async function getUserData() {
    await UserService.getUserData()
      .then(response =>
        user.value=response
      );
  }

  async function logout() {
    await AuthenticationService.logout()
      .then(() => {})
      .catch(() => {})
      .finally(async () => {
        user.value = null;
      });
  }
  const usersInitials = computed(() => {
    return `${user.value?.name[0]}${user.value?.surname[0]}`.toUpperCase()
  });

  return { user, getUserData, logout, usersInitials };
  },
  {
    persist: {
      key: 'roomplay-user',
  }
});

