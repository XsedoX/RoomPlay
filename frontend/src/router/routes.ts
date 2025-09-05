import type { RouteRecordRaw } from 'vue-router';
import LoginPage from '@/login_page/LoginPage.vue';

export const routes: RouteRecordRaw[] = [
  {path: '/', component: LoginPage}
]
