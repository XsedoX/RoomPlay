<script setup lang="ts">
import { LoginRepository } from '@/infrastructure/repositories/login_repository.ts';
import { ref } from 'vue'
import type IUserDataResponse from '@/infrastructure/models/IUserDataResponse.ts';

const userData = ref<IUserDataResponse | null>(null);

LoginRepository.userData()
  .then(data => { userData.value = data; })
  .catch(() => { userData.value = null; });
</script>

<template>
  <div v-if="userData">
    {{ userData.name }} {{ userData.surname }} {{ userData.role }} {{ userData.roomId }}
  </div>
  <div v-else>
    Authenticating...
  </div>
</template>

<style scoped></style>
