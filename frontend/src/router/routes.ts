import type { RouteRecordRaw } from 'vue-router';
import LoginPage from '@/pages/login_page/LoginPage.vue';
import MainMenuPage from '@/pages/main_menu_page/MainMenuPage.vue';
import RoomPage from '@/pages/room_page/RoomPage.vue';
import SettingsPage from '@/settings_page/SettingsPage.vue';
import SignInOidcPage from '@/pages/authentication_page/SignInOidcPage.vue'

export const routes: RouteRecordRaw[] = [
  {path: '/', component: LoginPage, name: 'LoginPage'},
  {path: '/mainMenu', component: MainMenuPage, name: 'MainMenuPage'},
  {
    path: '/room/:id',
    component: RoomPage,
    name: 'RoomPage'},
  {
    path: '/room/:id/settings',
    component: SettingsPage,
    name: 'SettingsPage'
  },
  {path: "/signin-oidc", component: SignInOidcPage,name: "SignInOidcPage" }
]
