<script setup lang="ts">
interface AvatarFullNameValueProps {
  fullName: string
  avatarAbbreviation: string
}

const props = defineProps<AvatarFullNameValueProps>();
</script>

<template>
  <v-container>
    <v-row justify="end" no-wrap>
      <v-col cols="auto">
        <div class="custom-avatar-card">
          <span class="text-body-1 text-nowrap text-black">{{ props.fullName }}</span>
          <v-avatar color="primary" class="custom-avatar pa-1">{{ props.avatarAbbreviation }}</v-avatar>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
  .custom-avatar-card {
    display: flex;
    align-items: center;
    padding-left: 8px;
    border: thin solid rgb(var(--v-theme-primary));
    border-radius: 24px;
    gap: 4px;
  }
  .custom-avatar {
    margin-right: -1px;
    margin-top: -1px;
    margin-bottom: -1px;
  }
</style>
