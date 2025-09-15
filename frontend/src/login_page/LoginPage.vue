<script setup lang="ts">
import LogoWithTitleText from '@/shared/LogoWithTitleText.vue';
</script>

<template>
  <v-container class="fill-height">
    <v-row align="center"
           justify="center">
      <v-col cols="8"
             sm="6"
             md="5"
             class="text-center">
        <LogoWithTitleText/>
        <v-btn prepend-icon="$googleIcon"
               class="text-body-1"
               rounded="xl"
               data-testid="login-with-google-btn"
               size="large">
          Continue with Google
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

