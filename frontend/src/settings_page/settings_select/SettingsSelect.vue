<script setup lang="ts">

import type ISettingsSelectProps from '@/settings_page/settings_select/ISettingsSelectProps.ts';
import type { TSettingsSelectEmits } from '@/settings_page/settings_select/TSettingsSelectEmits.ts';

const props = defineProps<ISettingsSelectProps>();
defineEmits<TSettingsSelectEmits>();
</script>

<template>
  <v-select
    rounded="xl"
    :hide-details="true"
    flat
    :model-value="props.value"
    @update:model-value="(newValue)=>$emit('on-selected', newValue)"
    :items="props.items"
    variant="solo"
    class="text-primary d-flex justify-end flex-shrink-0"
    bg-color="transparent"
  >
    <template v-slot:selection="{ item }">
      <span class="text-primary">{{ item.value }}</span>
    </template>
    <template v-slot:append-inner>
      <slot></slot>
    </template>
    <template v-slot:item="{ props }">
      <v-list-item v-bind="props" class="text-primary"> </v-list-item>
    </template>
  </v-select>
</template>
