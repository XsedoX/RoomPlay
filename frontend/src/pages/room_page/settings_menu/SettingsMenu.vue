<script setup lang="ts">
import PopupBase from '@/shared/popup_base/PopupBase.vue';
import { ref } from 'vue';
import { useTheme } from 'vuetify';
import { MenuItemsTypes } from '@/pages/room_page/settings_menu/MenuItemsTypes.ts'

const isAdmin = true;
function onMenuItemClick(id: MenuItemsTypes) {
  switch (id) {
    case MenuItemsTypes.QrCode:
      popup.value?.open();
      break;
    case MenuItemsTypes.Logout:
      console.log("Logout clicked");
      break;
    case MenuItemsTypes.ThemeMode:
      theme.toggle();
      break;
  }
}
const popup = ref<InstanceType<typeof PopupBase> | null>(null);
const theme = useTheme();
</script>

<template>
  <v-menu activator="parent">
    <v-list density="compact"
            @click:select="(value)=>onMenuItemClick(value.id as MenuItemsTypes)">
      <v-list-item :key="MenuItemsTypes.ThemeMode"
                   slim
                   :prepend-icon="theme.global.current.value.dark ? 'light_mode': 'dark_mode'"
                   :value="MenuItemsTypes.ThemeMode">
        {{theme.global.current.value.dark ? "Light Mode": "Dark Mode"}}
      </v-list-item>
      <v-list-item :key="MenuItemsTypes.QrCode"
                   slim
                   prepend-icon="qr_code"
                   :value="MenuItemsTypes.QrCode">
        QR Code
      </v-list-item>
      <v-list-item :key="MenuItemsTypes.Settings"
                   v-if="isAdmin"
                   slim
                   :to="`${$route.params['id']}/settings`"
                   prepend-icon="settings"
                   :value="MenuItemsTypes.Settings">
        Settings
      </v-list-item>
      <v-list-item :key="MenuItemsTypes.Logout"
                   to="/"
                   slim
                   prepend-icon="logout"
                   :value="MenuItemsTypes.Logout"
                   base-color="red">
        Logout
      </v-list-item>
    </v-list>
  </v-menu>
  <v-theme-provider theme="light">
    <popup-base ref="popup"
                popup-title="Scan the QR Code">
      <v-img cover
             rounded="xl"
             src="https://picsum.photos/200"
             alt="QR Code"
             aspect-ratio="1/1"
             width="240"/>
    </popup-base>
  </v-theme-provider>
</template>
<style scoped>
@import '@/assets/shared.css';
</style>

