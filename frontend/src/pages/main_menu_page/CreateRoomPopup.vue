<script setup lang="ts">
import { shallowRef } from 'vue';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as z from 'zod';
import type ICreateRoomRequest from '@/infrastructure/room/ICreateRoomRequest.ts';
import { useRoomStore } from '@/stores/room_store.ts';

const isPasswordVisible = shallowRef(false);
const isRepeatPasswordVisible = shallowRef(false);
const dialog = shallowRef(false);
const roomStore = useRoomStore();

const validationSchema = toTypedSchema(
  z
    .object({
      roomName: z
        .string()
        .min(5, 'Room name has to have at least 5 characters')
        .max(30, 'Room name has to have at most 30 characters'),
      roomPassword: z
        .string()
        .min(10, 'Password has to have at least 10 characters')
        .max(30, 'Password has to have at most 30 characters')
        .regex(/^\S*$/, 'Password cannot contain whitespaces'),
      repeatRoomPassword: z.string(),
    })
    .refine((data) => data.roomPassword === data.repeatRoomPassword, {
      message: "Passwords don't match",
      path: ['repeatRoomPassword'],
    }),
);
const { handleSubmit, defineField, setErrors, errors } = useForm({
  validationSchema,
  initialValues: {
    roomName: '',
    roomPassword: '',
    repeatRoomPassword: '',
  },
});

const [roomName] = defineField('roomName');
const [roomPassword] = defineField('roomPassword');
const [repeatRoomPassword] = defineField('repeatRoomPassword');

const onSubmit = handleSubmit(async (values) => {
  const request: ICreateRoomRequest = {
    roomName: values.roomName,
    roomPassword: values.roomPassword,
    repeatRoomPassword: values.repeatRoomPassword,
  };
  const validationErrors = await roomStore.createRoom(request);
  if (validationErrors) {
    setErrors(validationErrors);
    return;
  }
});
</script>

<template>
  <v-dialog
    data-testid="create-room-popup"
    activator="parent"
    max-width="290"
    v-model="dialog"
    persistent
  >
    <template v-slot:default>
      <v-card rounded="xl" class="pa-4">
        <v-container class="pa-0">
          <v-row justify="center" align="center" no-gutters>
            <v-col cols="2"></v-col>
            <v-col cols="8" class="text-center">
              <span data-testid="create-room-popup-title" class="text-h5 text-no-wrap"
                >Create a Room</span
              >
            </v-col>
            <v-col cols="2" class="d-flex justify-end">
              <v-btn icon="close" variant="text" size="small" @click="dialog = false"></v-btn>
            </v-col>
          </v-row>
          <v-row class="pt-1 pb-4" no-gutters
            ><v-col><v-divider></v-divider></v-col
          ></v-row>
        </v-container>
        <v-container class="pa-0">
          <v-form @submit.prevent="onSubmit">
            <v-row justify="center" no-gutters class="pb-4">
              <v-col>
                <v-text-field
                  v-model="roomName"
                  data-testid="create-room-popup-name-input"
                  :error-messages="errors.roomName ?? []"
                  label="Room Name"
                  required
                  hint="This cannot be changed later"
                  persistent-hint
                  clearable
                  clear-icon="close"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" no-gutters class="pb-6">
              <v-col>
                <v-text-field
                  v-model="roomPassword"
                  data-testid="create-room-popup-password-input"
                  :error-messages="errors.roomPassword ?? []"
                  label="Password"
                  hint="This cannot be changed later"
                  persistent-hint
                  required
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'visibility' : 'visibility_off'"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" no-gutters class="pb-6">
              <v-col>
                <v-text-field
                  label="Repeat Password"
                  data-testid="create-room-popup-repeat-password-input"
                  v-model="repeatRoomPassword"
                  :error-messages="errors.repeatRoomPassword ?? []"
                  required
                  :type="isRepeatPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isRepeatPasswordVisible ? 'visibility' : 'visibility_off'"
                  @click:append-inner="isRepeatPasswordVisible = !isRepeatPasswordVisible"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
              <v-col>
                <v-btn color="primary"
                       type="submit"
                       data-testid="create-room-popup-btn"
                       block
                       rounded="xl"> Create </v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-container>
      </v-card>
    </template>
  </v-dialog>
</template>
