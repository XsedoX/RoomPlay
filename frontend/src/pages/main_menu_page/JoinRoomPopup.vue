<script setup lang="ts">
import { shallowRef } from 'vue';
import { PlatformDiscoverer } from '@/infrastructure/utils/platform_discoverer.ts';
import { useForm } from 'vee-validate';
import { THostDevice } from '@/pages/settings_page/choose_host_device_list/THostDevice.ts';
import { useRoomStore } from '@/stores/room_store.ts';
import { toTypedSchema } from '@vee-validate/zod';
import * as z from 'zod';
import type IJoinRoomPasswordRequest from '@/infrastructure/room/IJoinRoomPasswordRequest';

const dialog = shallowRef(false);
const isPasswordVisible = shallowRef(false);
const roomStore = useRoomStore();

const validationSchema = toTypedSchema(
  z.object({
    roomName: z
      .string()
      .min(5, 'Room name has to have at least 5 characters')
      .max(30, 'Room name has to have at most 30 characters'),
    roomPassword: z
      .string()
      .min(10, 'Password has to have at least 10 characters')
      .max(30, 'Password has to have at most 30 characters')
      .regex(/^\S*$/, 'Password cannot contain whitespaces'),
  }),
);
const { handleSubmit, defineField, setErrors, errors } = useForm({
  validationSchema,
  initialValues: {
    roomName: '',
    roomPassword: '',
  },
});
const [roomName] = defineField('roomName');
const [roomPassword] = defineField('roomPassword');

const onSubmit = handleSubmit(async (values) => {
  const request: IJoinRoomPasswordRequest = {
    roomName: values.roomName,
    roomPassword: values.roomPassword,
  };
  const validationErrors = await roomStore.joinRoomPassword(request);
  if (validationErrors) {
    setErrors(validationErrors);
    return;
  }
});
</script>

<template>
  <v-dialog
    activator="parent"
    data-testid="join-room-dialog"
    max-width="290"
    v-model="dialog"
    persistent
  >
    <template v-slot:default>
      <v-card rounded="xl" class="pa-4">
        <v-container class="pa-0">
          <v-row justify="center" align="center" no-gutters>
            <v-col cols="2"></v-col>
            <v-col cols="8" data-testid="join-room-dialog-title" class="text-center">
              <span class="text-h5">Join a Room</span>
            </v-col>
            <v-col cols="2" class="d-flex justify-end">
              <v-btn icon="close" variant="text" size="small" @click="dialog = false"></v-btn>
            </v-col>
          </v-row>
          <v-row class="pt-1 pb-4" no-gutters>
            <v-col><v-divider></v-divider></v-col>
          </v-row>
        </v-container>
        <v-container class="pa-0">
          <v-form @submit.prevent="onSubmit">
            <v-row justify="center" no-gutters>
              <v-col>
                <v-text-field
                  v-model="roomName"
                  :error-messages="errors.roomName ?? []"
                  label="Room Name"
                  data-testid="join-room-popup-name-input"
                  required
                  clearable
                  clear-icon="close"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
              <v-col>
                <v-text-field
                  v-model="roomPassword"
                  :error-messages="errors.roomPassword ?? []"
                  label="Password"
                  data-testid="join-room-popup-password-input"
                  required
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'visibility' : 'visibility_off'"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                ></v-text-field>
              </v-col>
            </v-row>
            <div v-if="PlatformDiscoverer.getDeviceType() === THostDevice.Mobile">
              <v-row justify="center" no-gutters class="pb-2">
                <v-col>
                  <v-divider>
                    <span class="text-body-1">or</span>
                  </v-divider>
                </v-col>
              </v-row>
              <v-row justify="center" no-gutters>
                <v-col>
                  <v-btn
                    variant="plain"
                    color="primary"
                    block
                    :ripple="false"
                    prepend-icon="qr_code"
                  >
                    Scan a QR Code
                  </v-btn>
                </v-col>
              </v-row>
            </div>
            <v-row justify="center">
              <v-col>
                <v-btn
                  data-testid="join-room-popup-btn"
                  type="submit"
                  color="primary"
                  block
                  rounded="xl"
                >
                  Join
                </v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-container>
      </v-card>
    </template>
  </v-dialog>
</template>
