<script setup lang="ts">
import { shallowRef } from 'vue';
import type ITouchscreenTooltipProps from '@/shared/TouchscreenTooltip/ITouchscreenTooltipProps.ts';

const showTooltip = shallowRef(false);
const props = withDefaults(defineProps<ITouchscreenTooltipProps>(),{
  openOnHover: true
});
</script>

<template>
  <v-tooltip
    v-model="showTooltip"
    @click:outside="showTooltip = false"
    :open-on-hover="props.openOnHover"
    open-on-click
    :text="props.text"
  >
    <template v-slot:activator="{ props: tooltipProps }">
      <slot :tooltipProps='tooltipProps'></slot>
    </template>
  </v-tooltip>
</template>

