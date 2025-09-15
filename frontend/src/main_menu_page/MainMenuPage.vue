<script setup lang="ts">
import LogoWithTitleText from '@/shared/LogoWithTitleText.vue';
import AvatarWithFullName from '@/main_menu_page/AvatarWithFullName.vue';
const fullName = "Full Name";
const avatarAbbreviation = "FN";
</script>

<template>
  <v-container fluid>
    <v-row align="end"
           justify="end"
           dense>
      <v-col cols="auto"
             class="text-center">
        <AvatarWithFullName :full-name="fullName" :avatar-abbreviation="avatarAbbreviation"></AvatarWithFullName>
      </v-col>
    </v-row>
    <v-row align="center"
           justify="center"
           dense>
      <v-col cols="8"
             sm="6"
             md="3"
             class="text-center pa-2">
        <LogoWithTitleText></LogoWithTitleText>
      </v-col>
    </v-row>
    <v-row align="center"
           justify="center"
           dense>
      <v-col cols="8"
             sm="6"
             md="3"
             class="text-center">
        <v-btn data-testid = "join-room-btn"
               rounded="xl"
               size="x-large"
               color="primary">
          Join a Room
        </v-btn>
      </v-col>
    </v-row>
    <v-row align="center"
           justify="center"
    >
      <v-col cols="8"
             sm="6"
             md="3"
             class="text-center">
        <v-btn data-testid = "create-room-btn"
               rounded="xl"
               variant="outlined"
               size="x-large"
               color="primary">
          Create a Room
        </v-btn>
      </v-col>
    </v-row>
    <v-row align="center"
           justify="center"
           dense>
      <v-col cols="8"
             sm="6"
             md="3"
             class="text-center">
        <v-btn data-testid = "logout-btn"
               variant="text"
               rounded="xl"
               size="large"
               color="error">
          Logout
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

