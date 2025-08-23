import type { RouteRecordRaw } from 'vue-router';
import LoginPage from '@/loginPage/LoginPage.vue';

export const routes: RouteRecordRaw[] = [
  {path: '/', component: LoginPage}
]
