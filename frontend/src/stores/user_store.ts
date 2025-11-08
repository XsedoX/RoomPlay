import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type { ILoggedInUserStoreModel } from '@/infrastructure/models/ILoggedInUserStoreModel.ts';
import { UserService } from '@/infrastructure/services/user_service.ts';
import { LoginService } from '@/infrastructure/services/login_service.ts';
import { useRouter } from 'vue-router'

export const useUserStore = defineStore('user',  () => {
    const user = ref<ILoggedInUserStoreModel | null>(null);
    const router = useRouter();

    async function getUserData() {
      await UserService.getUserData()
        .then(response =>
          user.value=response
        );
    }

    async function logout() {
      await LoginService.logout()
        .then(() => {})
        .catch(() => {})
        .finally(async () => {
          user.value = null;
          await router.replace({ name: 'LoginPage' });});
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

