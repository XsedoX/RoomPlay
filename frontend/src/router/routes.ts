import type { RouteRecordRaw } from 'vue-router';
import LoginPage from '@/login_page/LoginPage.vue';
import MainMenuPage from '@/main_menu_page/MainMenuPage.vue';

export const routes: RouteRecordRaw[] = [
  {path: '/', component: LoginPage},
  {path: '/mainMenu', component: MainMenuPage}
]
