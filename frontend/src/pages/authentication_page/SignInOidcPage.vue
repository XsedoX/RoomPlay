<script setup lang="ts">
import { useUserStore } from '@/stores/user_store.ts';
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const userStore = useUserStore();

onMounted(async () => {
  await userStore.getUserData()
  await router.replace({ name: 'MainMenuPage' })
})
</script>

<template>
  Authenticating...
</template>

<style scoped></style>
