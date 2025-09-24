<script setup lang="ts">
import { useTheme } from 'vuetify'

interface AvatarFullNameValueProps {
  fullName: string
  avatarAbbreviation: string
}

const props = defineProps<AvatarFullNameValueProps>();
const theme = useTheme()
const primaryColor = theme.current.value.colors.primary
</script>

<template>
  <v-container class="pa-0 ma-0">
    <v-row justify="end" no-wrap no-gutters>
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
    border: thin solid v-bind(primaryColor);
    border-radius: 24px;
    gap: 4px;
  }
  .custom-avatar {
    margin-right: -1px;
    margin-top: -1px;
    margin-bottom: -1px;
  }
</style>
