<script setup lang="ts">
import { useNotificationStore } from '@/stores/notification_store.ts';

const notificationStore = useNotificationStore();
</script>

<template>
  <v-app min-width="320px">
    <RouterView />
  </v-app>
  <v-snackbar
    v-model="notificationStore.snackbarVisible"
    :color="notificationStore.snackbarColor"
    :timeout="3000"
    location="top right"
  >
    {{ notificationStore.snackbarMessage }}
  </v-snackbar>
</template>
