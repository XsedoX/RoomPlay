import { createRouter, createWebHistory } from 'vue-router';
import { routes } from '@/router/routes.ts';
import { useUserStore } from '@/stores/user_store.ts';
import { useRoomStore } from '@/stores/room_store.ts';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
});

router.beforeEach((to, _, next) => {
  const userStore = useUserStore();
  const roomStore = useRoomStore();
  const forbiddenRoutesForAuthUser = ['LoginPage', 'SignInOidcPage'];
  const forbiddenRoutesForAuthenticatedUserInRoom = ['MainMenuPage'];

  // Redirect to the main menu if the user is authenticated and tries to access the login page
  if (to.meta['requiresAuth'] && !userStore.user) {
    next({ name: 'LoginPage' });
  } else if (userStore.user && to.name && forbiddenRoutesForAuthenticatedUserInRoom.includes(to.name.toString()) && roomStore.room !== null) {
    next({ name: 'RoomPage'});
  }
  else if (userStore.user && to.name && forbiddenRoutesForAuthUser.includes(to.name.toString()) && roomStore.room === null) {
    next({ name: 'MainMenuPage' });
  } else {
    next();
  }
});

export default router;
