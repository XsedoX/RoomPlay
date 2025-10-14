import type { RouteRecordRaw } from 'vue-router';
import LoginPage from '@/login_page/LoginPage.vue';
import MainMenuPage from '@/main_menu_page/MainMenuPage.vue';
import RoomPage from '@/room_page/RoomPage.vue';
import SettingsPage from '@/settings_page/SettingsPage.vue';

export const routes: RouteRecordRaw[] = [
  {path: '/', component: LoginPage, name: 'login'},
  {path: '/mainMenu', component: MainMenuPage, name: 'mainMenu'},
  {
    path: '/room/:id',
    component: RoomPage,
    name: 'room'},
  {
    path: '/room/:id/settings',
    component: SettingsPage,
    name: 'roomSettings'
  }
]
