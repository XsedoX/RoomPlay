<script setup lang="ts">
import MenuItems from '@/room_page/settings_menu/MenuItems.ts';
import { shallowRef } from 'vue';

const showQrCodeDialog = shallowRef(false);
const isAdmin = true;
function onMenuItemClick(id: unknown) {
  switch (id) {
    case MenuItems.QR_CODE:
      showQrCodeDialog.value = true;
      break;
    case MenuItems.LOGOUT:
      console.log("Logout clicked");
      break;
  }
}
</script>

<template>
  <v-menu activator="parent">
    <v-list density="compact"
            @click:select="(value)=>onMenuItemClick(value.id)">
      <v-list-item :key="MenuItems.QR_CODE"
                   slim
                   prepend-icon="qr_code"
                   :value="MenuItems.QR_CODE">
        QR Code
      </v-list-item>
      <v-list-item :key="MenuItems.SETTINGS"
                   v-if="isAdmin"
                   slim
                   :to="`${$route.params['id']}/settings`"
                   prepend-icon="settings"
                   :value="MenuItems.SETTINGS">
        Settings
      </v-list-item>
      <v-list-item :key="MenuItems.LOGOUT"
                   to="/"
                   slim
                   prepend-icon="logout"
                   :value="MenuItems.LOGOUT"
                   base-color="red">
        Logout
      </v-list-item>
    </v-list>
  </v-menu>
  <v-dialog v-model="showQrCodeDialog" max-width="300">
    <v-card rounded="xl">
      <v-card-title class="d-flex justify-space-between align-center px-1 pb-0 pt-1">
        <v-btn icon="close" variant="text" :disabled="true" class="invisible"></v-btn>
        <div class="text-center">Scan the QR Code</div>
        <v-btn icon="close" variant="text" @click="showQrCodeDialog = false"></v-btn>
      </v-card-title>
      <v-divider class="mx-2"></v-divider>
      <v-card-text class="d-flex justify-center align-center pa-4" >
        <v-img cover rounded="xl" src="https://picsum.photos/200" alt="QR Code" aspect-ratio="1/1" width="240"/>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>
<style scoped>
@import '@/assets/shared.css';
</style>

