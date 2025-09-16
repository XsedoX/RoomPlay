<script setup lang="ts">
import { shallowRef } from 'vue'

const dialog = shallowRef(false)
const isPasswordVisible = shallowRef(false)
</script>

<template>
  <v-dialog activator="parent" max-width="290" v-model="dialog" persistent>
    <template v-slot:default>
      <v-card>
        <v-row justify="center" align="center" no-gutters>
          <v-col cols="2"></v-col>
          <v-col cols="8" class="text-center">
            <span class="text-h5">Join a Room</span>
          </v-col>
          <v-col cols="2" class="d-flex justify-end">
            <v-btn icon="close" variant="text" size="small" @click="dialog = false"></v-btn>
          </v-col>
        </v-row>
        <v-row class="pb-4" no-gutters><v-col><v-divider></v-divider></v-col></v-row>
        <v-form>
          <v-row justify="center" no-gutters>
            <v-col class="px-6">
              <v-text-field
                label="Room Name"
                required
                clearable
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row justify="center" no-gutters>
            <v-col class="px-6">
              <v-text-field
                label="Password"
                required
                :type="isPasswordVisible ? 'text' : 'password'"
                :append-inner-icon="isPasswordVisible ? 'visibility' : 'visibility_off'"
                @click:append-inner="isPasswordVisible = !isPasswordVisible"
              ></v-text-field>
            </v-col>
          </v-row>
        </v-form>
      </v-card>
    </template>
  </v-dialog>
</template>

<style scoped></style>
