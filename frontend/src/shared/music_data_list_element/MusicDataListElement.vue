<script setup lang="ts">

import type { IMusicDataListElementProps } from '@/shared/music_data_list_element/IMusicDataListElementProps.ts';

const props = defineProps<IMusicDataListElementProps>()
</script>
<template>
  <v-sheet color="surface-container"
           @click="(event:Event)=> props.onElementClick({event: event, id: props.musicDataListElementViewModel.id})"
           class="d-flex justify-start align-center pa-3 ga-2">
    <div class="overflow-hidden">
      <v-img
        cover
        max-width="75px"
        min-width="50px"
        rounded="lg"
        aspect-ratio="1/1"
        :src="props.musicDataListElementViewModel.imageUrl">
      </v-img>
    </div>
    <div class="d-flex flex-column justify-between">
      <span class="text-truncate">
            {{ props.musicDataListElementViewModel.title }}
    </span>
      <span class="on-surface-variant text-truncate">
      {{
          props.musicDataListElementViewModel.subtitle
        }}
    </span>
    </div>
  </v-sheet>
</template>

