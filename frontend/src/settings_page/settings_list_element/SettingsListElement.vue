<script setup lang="ts">

import type ISettingsListElementProps from '@/settings_page/settings_list_element/ISettingsListElementProps.ts';

const props = defineProps<ISettingsListElementProps>()
</script>

<template>
  <v-list-subheader>{{ props.header }}</v-list-subheader>
  <v-list-item>
    <v-sheet rounded="xl"
             color="surface-container pa-2 d-flex flex-nowrap ga-2 align-center">
      <div class="d-flex flex-column mr-auto">
        <div class="h6 font-weight-bold">{{ props.subHeader }}</div>
        <div v-if="props.hint !== null"
             class="text-medium-emphasis text-caption">{{props.hint}}</div>
      </div>
      <slot></slot>
    </v-sheet>
  </v-list-item>
</template>
