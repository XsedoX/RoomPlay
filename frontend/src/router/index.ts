import { createRouter, createWebHistory } from 'vue-router';
import { routes } from '@/router/routes.ts';
import { useUserStore } from '@/stores/user_store.ts';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes
});

router.beforeEach((to, _, next) => {
  const userStore = useUserStore();
  const forbiddenRoutesForAuthUser = ['LoginPage', 'SignInOidcPage'];

  if (userStore.user && to.name && forbiddenRoutesForAuthUser.includes(to.name.toString())) {
    next({ name: 'MainMenuPage' });
  } else if (to.meta['requiresAuth'] && !userStore.user) {
    next({ name: 'LoginPage' });
  } else {
    next();
  }
});

export default router;
