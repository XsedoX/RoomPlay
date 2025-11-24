<script setup lang="ts">
import TouchscreenTooltip from '@/shared/touchscreen_tooltip/TouchscreenTooltip.vue';
import type IPageTitleProps from '@/shared/page_title/IPageTitleProps.ts';

const props = defineProps<IPageTitleProps>();
</script>

<template>
  <v-sheet
    color="surface-container-highest"
    class="elevation-4 pa-1 d-flex align-center justify-space-between ga-1"
  >
    <slot name="top-left-corner"></slot>
    <touchscreen-tooltip :open-on-hover="false" :text="props.title" v-slot="{ tooltipProps }">
      <div
        data-testid="page-title"
        v-bind="tooltipProps"
        class="text-h6 w-0 text-no-wrap text-truncate text-center flex-grow-1"
      >
        {{ props.title }}
      </div>
    </touchscreen-tooltip>
    <slot name="top-right-corner"></slot>
  </v-sheet>
</template>
