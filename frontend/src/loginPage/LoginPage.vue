<script setup lang="ts">
import logo from '@/assets/logo.png'
</script>

<template>
  <v-container class="fill-height">
      <v-row align="center" justify="center">
        <v-col
          cols="8"
          sm="6"
          md="4"
          xl="3"
        >
          <v-img
            :src="logo"
            alt="RoomPlay2 logo"
            width="100%"
            aspect-ratio="1"></v-img>
        </v-col>
      </v-row>
  </v-container>
</template>

<style scoped></style>
