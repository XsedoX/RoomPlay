<script setup lang="ts">

import type ISettingsListElementProps from '@/settings_page/settings_list_element/ISettingsListElementProps.ts';

const props = defineProps<ISettingsListElementProps>();
</script>

<template>
  <v-list-subheader>{{ props.header }}</v-list-subheader>
  <v-list-item>
    <v-sheet rounded="lg"
             color="surface-container justify-start pa-2 d-flex flex-nowrap ga-2 align-center">
      <slot name="content">
        <div class="d-flex align-center ga-2 mr-auto justify-start w-75">
          <slot name="image"></slot>
          <div class="d-flex flex-column flex-grow-1 w-0">
            <div :class="['h6', 'font-weight-bold']">
              {{ props.subHeader }}
            </div>
            <div :class="['text-medium-emphasis', 'text-caption']">
              {{props.hint}}
            </div>
          </div>
        </div>
        <slot></slot>
      </slot>
    </v-sheet>
  </v-list-item>
</template>
