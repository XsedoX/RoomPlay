<script setup lang="ts">
import { shallowRef } from 'vue';

const isPasswordVisible = shallowRef(false)
const dialog = shallowRef(false)
</script>

<template>
  <v-dialog activator="parent" max-width="290" v-model="dialog" persistent>
    <template v-slot:default>
      <v-card rounded="xl" class="pa-4">
        <v-container class="pa-0">
          <v-row justify="center"
                 align="center"
                 no-gutters>
            <v-col cols="2"></v-col>
            <v-col cols="8" class="text-center">
              <span class="text-h5 text-no-wrap">Create a Room</span>
            </v-col>
            <v-col cols="2" class="d-flex justify-end">
              <v-btn icon="close" variant="text" size="small" @click="dialog = false"></v-btn>
            </v-col>
          </v-row>
          <v-row class="pt-1 pb-4" no-gutters><v-col><v-divider></v-divider></v-col></v-row>
        </v-container>
        <v-container class="pa-0">
          <v-form>
            <v-row justify="center"
                   no-gutters
                   class="pb-4">
              <v-col>
                <v-text-field
                  label="Room Name"
                  required
                  hint="This cannot be changed later"
                  persistent-hint
                  clearable
                  clear-icon="close"></v-text-field>
              </v-col>
            </v-row>
            <v-row justify="center"
                   no-gutters
                   class="pb-6">
              <v-col>
                <v-text-field
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
            <v-row justify="center" no-gutters>
              <v-col>
                <v-btn color="primary"
                       block
                       rounded="xl">
                  Create
                </v-btn>
              </v-col>
            </v-row>
          </v-form>
        </v-container>
      </v-card>
    </template>
  </v-dialog>
</template>
