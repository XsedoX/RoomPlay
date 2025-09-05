<script setup lang="ts">
import logo from '@/assets/logo.png'
</script>

<template>
  <v-container class="fill-height">
    <v-row align="center"
      justify="center">
      <v-col
        cols="8"
        sm="6"
        md="5"
        class="text-center"
      >
        <v-img
          :src="logo"
          alt="RoomPlay2 logo"
          data-testid="logo"
          aspect-ratio="1"
          class="w-100 mb-6"/>
        <p class="text-h3">RoomPlay2</p>
        <p class="text-subtitle-1 mb-6">The playlist that creates itself.</p>
        <v-btn prepend-icon="$googleIcon"
               class="pa-2 text-body-1"
               data-testid="login-with-google-btn"
               size="medium">
          Continue with Google
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped></style>
