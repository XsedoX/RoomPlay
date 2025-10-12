<script setup lang="ts">
import { shallowRef } from 'vue';
import type IPopupBaseProps from '@/shared/popup_base/IPopupBaseProps.ts';

const showQrCodeDialog = shallowRef(false);
function open() {
  showQrCodeDialog.value = true;
}
const props = defineProps<IPopupBaseProps>();
defineExpose({
  open
});
</script>

<template>
  <v-dialog v-model="showQrCodeDialog"
            max-width="300">
    <v-card rounded="xl">
      <v-card-title class="d-flex justify-space-between align-center px-1 pb-0 pt-1">
        <v-btn icon="close"
               variant="text"
               :disabled="true"
               class="invisible"></v-btn>
        <div class="text-center">
          {{ props.popupTitle }}
        </div>
        <v-btn icon="close"
               variant="text"
               @click="showQrCodeDialog = false"></v-btn>
      </v-card-title>
      <v-divider class="mx-2"></v-divider>
      <v-card-text class="d-flex justify-center align-center pa-4">
        <slot></slot>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<style scoped>
@import '@/assets/shared.css';
</style>
