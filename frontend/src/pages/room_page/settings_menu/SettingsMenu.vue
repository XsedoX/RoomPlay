<script setup lang="ts">
import PopupBase from '@/shared/popup_base/PopupBase.vue';
import { computed, ref } from 'vue';
import { useTheme } from 'vuetify';
import { MenuItemsTypes } from '@/pages/room_page/settings_menu/MenuItemsTypes.ts';
import { useRoomStore } from '@/stores/room_store.ts';
import { useQRCode } from '@vueuse/integrations/useQRCode';
import { useUserStore } from '@/stores/user_store.ts';
import { useRouter } from 'vue-router';

const roomStore = useRoomStore();
const userStore = useUserStore();
const router = useRouter();
const qrCodeData = computed(() => (roomStore.room?.qrCode ? roomStore.room.qrCode : ''));
const qrCode = useQRCode(qrCodeData);

async function onMenuItemClick(id: MenuItemsTypes) {
  switch (id) {
    case MenuItemsTypes.QrCode:
      popup.value?.open();
      break;
    case MenuItemsTypes.Logout:
      await logout();
      break;
    case MenuItemsTypes.ThemeMode:
      theme.toggle();
      break;
  }
}
async function logout() {
  await userStore.logout();
  await router.replace({ name: 'LoginPage' });
}
const popup = ref<InstanceType<typeof PopupBase> | null>(null);
const theme = useTheme();
</script>

<template>
  <v-menu activator="parent">
    <v-list
      density="compact"
      @click:select="(value) => onMenuItemClick(value.id as MenuItemsTypes)"
    >
      <v-list-item
        :key="MenuItemsTypes.ThemeMode"
        slim
        :prepend-icon="theme.global.current.value.dark ? 'light_mode' : 'dark_mode'"
        :value="MenuItemsTypes.ThemeMode"
      >
        {{ theme.global.current.value.dark ? 'Light Mode' : 'Dark Mode' }}
      </v-list-item>
      <v-list-item
        :key="MenuItemsTypes.QrCode"
        slim
        prepend-icon="qr_code"
        :value="MenuItemsTypes.QrCode"
      >
        QR Code
      </v-list-item>
      <v-list-item
        :key="MenuItemsTypes.Settings"
        v-if="roomStore.isHost"
        slim
        :to="`/room/settings`"
        prepend-icon="settings"
        :value="MenuItemsTypes.Settings"
      >
        Settings
      </v-list-item>
      <v-list-item
        :key="MenuItemsTypes.Logout"
        slim
        prepend-icon="logout"
        :value="MenuItemsTypes.Logout"
        base-color="red"
      >
        Logout
      </v-list-item>
    </v-list>
  </v-menu>
  <v-theme-provider theme="light">
    <popup-base ref="popup" popup-title="Scan the QR Code">
      <v-img cover rounded="xl" :src="qrCode" alt="QR Code" aspect-ratio="1/1" width="240" />
    </popup-base>
  </v-theme-provider>
</template>
<style scoped>
@import '@/assets/shared.css';
</style>
