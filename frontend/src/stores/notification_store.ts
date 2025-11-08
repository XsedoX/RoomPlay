import { defineStore } from 'pinia';
import { ref } from 'vue';
import { TSnackbarColor } from '@/infrastructure/models/TSnackbarColor.ts';

export const useNotificationStore = defineStore("notification", ()=>{
  const snackbarVisible = ref(false);
  const snackbarMessage = ref("");
  const snackbarColor = ref(TSnackbarColor.ERROR);

  function showSnackbar(message: string, color: TSnackbarColor) {
    if (snackbarVisible.value) {return;}
    snackbarMessage.value = message;
    snackbarColor.value = color;
    snackbarVisible.value = true;
  }

  return { snackbarVisible, snackbarMessage, snackbarColor, showSnackbar };
})
