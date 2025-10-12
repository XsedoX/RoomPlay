<script setup lang="ts">
import type { TMusicDataListElementEmits } from '@/shared/music_data_list_element/TMusicDataListElementEmits.ts';
import type IGuidEvent from '@/shared/IGuidEvent.ts';
import type IMusicDataListElementProps from '@/shared/music_data_list_element/IMusicDataListElementProps.ts';

const props = defineProps<IMusicDataListElementProps>();
const emit = defineEmits<TMusicDataListElementEmits>();
function chooseMusic(event: IGuidEvent) {
  emit('on-music-choice', event);
}
</script>

<template>
  <v-sheet @click="(event: Event)=>chooseMusic({event: event, id: props.musicDataListDto.id})"
           rounded="lg"
           color="surface-container justify-start pa-2 d-flex flex-nowrap align-center">
    <div class="d-flex align-center ga-2 justify-start w-100">
      <div class="overflow-hidden">
        <v-img
          cover
          max-width="75px"
          min-width="50px"
          rounded="lg"
          aspect-ratio="1/1"
          :src="props.musicDataListDto.imageUrl">
        </v-img>
      </div>
      <div class="d-flex flex-column flex-grow-1 w-0">
        <div :class="['h6', 'font-weight-bold']">
          {{ props.musicDataListDto.title }}
        </div>
        <div :class="['text-medium-emphasis', 'text-caption']">
          {{props.musicDataListDto.subtitle}}
        </div>
      </div>
      <v-icon color="outline" icon="chevron_right"></v-icon>
    </div>
  </v-sheet>
</template>
