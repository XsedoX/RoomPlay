<script setup lang="ts">
import logo from '@/assets/logo.png'
</script>

<template>
  <v-row no-gutters
         justify="center">
    <v-col>
      <v-img
        :src="logo"
        alt="RoomPlay2 logo"
        data-testid="logo"
        aspect-ratio="1"
        class="w-100 mb-6"/>
    </v-col>
  </v-row>
  <v-row no-gutters
         justify="center">
    <v-col cols="auto" class="text-center d-flex justify-center">
      <span class="text-h3">RoomPlay2</span>
    </v-col>
  </v-row>
  <v-row no-gutters
         justify="center">
    <v-col cols="auto" class="text-center">
      <span class="text-subtitle-1 on-surface-variant">The playlist that creates itself.</span>
    </v-col>
  </v-row>
</template>
