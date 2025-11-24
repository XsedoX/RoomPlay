import { defineStore } from 'pinia';
import { ref } from 'vue';
import { TSnackbarColor } from '@/infrastructure/utils/TSnackbarColor.ts';

export const useNotificationStore = defineStore('notification', () => {
  const snackbarVisible = ref(false);
  const snackbarMessage = ref('');
  const snackbarColor = ref(TSnackbarColor.ERROR);

  function showSnackbar(message: string | null | undefined, color: TSnackbarColor) {
    if (snackbarVisible.value) {
      return;
    }
    snackbarMessage.value = message ?? "An unexpected error occurred.";
    snackbarColor.value = color;
    snackbarVisible.value = true;
  }

  return { snackbarVisible, snackbarMessage, snackbarColor, showSnackbar };
});
