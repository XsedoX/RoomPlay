<script setup lang="ts">
import LogoWithTitleText from '@/shared/LogoWithTitleText.vue';
import { AuthenticationService } from '@/infrastructure/authentication/authentication_service.ts';
import { useUserStore } from '@/stores/user_store.ts';
import { useRouter } from 'vue-router';
const userStore = useUserStore();
const router = useRouter();

async function login() {
  if (!userStore.user) {
    const redirectUri = await AuthenticationService.loginWithGoogle();
    if (redirectUri) {
      globalThis.location.assign(redirectUri);
    }
  } else {
    await router.replace({ name: 'MainMenuPage' });
  }
}
</script>

<template>
  <v-container class="fill-height">
    <v-row>
      <v-col class="d-flex flex-column ga-6">
        <v-row justify="center" no-gutters>
          <v-col cols="8" sm="6" md="5">
            <LogoWithTitleText />
          </v-col>
        </v-row>
        <v-row justify="center" no-gutters>
          <v-col cols="8" sm="6" md="5" class="d-flex align-center justify-center">
            <v-btn
              prepend-icon="$googleIcon"
              class="text-body-1"
              rounded="xl"
              @click="login"
              min-width="220"
              data-testid="login-with-google-btn"
              size="large"
            >
              Continue with Google
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>
