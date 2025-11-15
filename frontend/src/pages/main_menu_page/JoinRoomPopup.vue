<script setup lang="ts">
import { shallowRef } from 'vue';

const dialog = shallowRef(false);
const isPasswordVisible = shallowRef(false);
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
            <v-col cols="8" class="text-center">
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
          <v-form>
            <v-row justify="center" no-gutters>
              <v-col>
                <v-text-field
                  label="Room Name"
                  required
                  clearable
                  clear-icon="close"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
              <v-col>
                <v-text-field
                  label="Password"
                  required
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'visibility' : 'visibility_off'"
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                ></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center" no-gutters class="pb-2">
              <v-col>
                <v-divider>
                  <span class="text-body-1">or</span>
                </v-divider>
              </v-col>
            </v-row>
            <v-row justify="center" no-gutters>
              <v-col>
                <v-btn variant="plain" color="primary" block :ripple="false" prepend-icon="qr_code">
                  Scan a QR Code
                </v-btn>
              </v-col>
            </v-row>
            <v-row justify="center">
              <v-col>
                <v-btn color="primary" block rounded="xl"> Join </v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-container>
      </v-card>
    </template>
  </v-dialog>
</template>
