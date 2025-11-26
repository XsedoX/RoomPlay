<script setup lang="ts">
import PopupBase from '@/shared/popup_base/PopupBase.vue';
import { computed, ref } from 'vue';
import { TMenuItems } from '@/pages/room_page/settings_menu/TMenuItems';
import { TThemes } from '@/pages/room_page/settings_menu/TThemes';
import { useRoomStore } from '@/stores/room_store.ts';
import { useQRCode } from '@vueuse/integrations/useQRCode';
import { useUserStore } from '@/stores/user_store.ts';
import { useRouter } from 'vue-router';
import { useThemeStore } from '@/stores/theme_store.ts';

const roomStore = useRoomStore();
const userStore = useUserStore();
const router = useRouter();
const qrCodeData = computed(() => (roomStore.room?.qrCode ? roomStore.room.qrCode : ''));
const qrCode = useQRCode(qrCodeData);
const themeStore = useThemeStore();

async function onMenuItemClick(id: TMenuItems) {
  switch (id) {
    case TMenuItems.QrCode:
      popup.value?.open();
      break;
    case TMenuItems.Logout:
      await logout();
      break;
    case TMenuItems.ThemeMode:
      themeStore.toggleTheme();
      break;
  }
}
async function logout() {
  await userStore.logout();
  await router.replace({ name: 'LoginPage' });
}
const popup = ref<InstanceType<typeof PopupBase> | null>(null);
</script>

<template>
  <v-menu activator="parent">
    <v-list density="compact" @click:select="(value) => onMenuItemClick(value.id as TMenuItems)">
      <v-list-item
        :key="TMenuItems.ThemeMode"
        slim
        :prepend-icon="themeStore.storedTheme === TThemes.DarkMode ? 'light_mode' : 'dark_mode'"
        :value="TMenuItems.ThemeMode"
      >
        {{ themeStore.storedTheme === TThemes.DarkMode ? 'Light Mode' : 'Dark Mode' }}
      </v-list-item>
      <v-list-item
        :key="TMenuItems.QrCode"
        slim
        data-testid="qr-code-menu-item"
        prepend-icon="qr_code"
        :value="TMenuItems.QrCode"
      >
        QR Code
      </v-list-item>
      <v-list-item
        :key="TMenuItems.Settings"
        v-if="roomStore.isHost"
        slim
        :to="`/room/settings`"
        prepend-icon="settings"
        :value="TMenuItems.Settings"
      >
        Settings
      </v-list-item>
      <v-list-item
        :key="TMenuItems.Logout"
        slim
        data-testid="logout-menu-item"
        prepend-icon="logout"
        :value="TMenuItems.Logout"
        base-color="red"
      >
        Logout
      </v-list-item>
    </v-list>
  </v-menu>
  <popup-base ref="popup" popup-title="Scan the QR Code">
    <v-img
      data-testid="qr-code-img"
      cover
      rounded="xl"
      :src="qrCode"
      alt="QR Code"
      aspect-ratio="1/1"
      width="240"
    />
  </popup-base>
</template>
<style scoped>
@import '@/assets/shared.css';
</style>
