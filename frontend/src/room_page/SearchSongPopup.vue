<script setup lang="ts">
import { shallowRef } from 'vue';

const dialog = shallowRef(false)
</script>

<template>
  <v-row>
    <v-col cols="11">
      <v-dialog activator="parent" v-model="dialog">
        <template v-slot:default>
          <v-card class="d-flex flex-column" rounded="xl">
            <div class="d-flex align-center justify-center">
              <span>Search for a Song</span>
              <v-btn variant="text" icon="close"></v-btn>
            </div>
          </v-card>
        </template>
      </v-dialog>
    </v-col>
  </v-row>
</template>

