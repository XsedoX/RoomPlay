<script setup lang="ts">
import type { ISongListElementProps } from '@/room_page/song_list_element/ISongListElementProps.ts';
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import TouchscreenTooltip from '@/shared/touchscreen_tooltip/TouchscreenTooltip.vue';

const props = defineProps<ISongListElementProps>();
const backgroundColor = computed(() =>
  props.songListDto.wasPlayed ? 'outline-variant' : 'surface-container',
);
const theme = useTheme();
const outlineColor = theme.current.value.colors['outline'];
const isAdmin = true
</script>

<template>
  <v-sheet class="d-flex px-4 py-2 justify-start align-center w-100" :color="backgroundColor">
    <div class="overflow-hidden">
      <v-img
        cover
        max-width="75px"
        min-width="50px"
        rounded="lg"
        aspect-ratio="1/1"
        :src="props.songListDto.albumCoverUrl"
      >
      </v-img>
    </div>
    <v-container
      class="w-0 flex-shrink-1 flex-grow-1 py-0 ma-0 mr-auto d-flex flex-column justify-between px-2"
    >
      <touchscreen-tooltip :open-on-hover="false"
                           :text="props.songListDto.title"
                           v-slot="{ tooltipProps }" >
        <span v-bind="tooltipProps" class="text-truncate">
            {{ props.songListDto.title }}</span
        >
      </touchscreen-tooltip>
      <touchscreen-tooltip :text="props.songListDto.author"
                           :open-on-hover="false"
                           v-slot="{ tooltipProps }">
        <span v-bind="tooltipProps" class="on-surface-variant text-truncate">{{
          props.songListDto.author
        }}</span>
      </touchscreen-tooltip>
    </v-container>
    <div class="d-flex flex-column justify-space-between align-end">
      <v-tooltip close-delay="2" open-on-click :text="props.songListDto.addedBy">
        <template v-slot:activator="{ props: tooltipProps }">
          <v-avatar v-bind="tooltipProps" size="x-small" color="primary">{{
            props.songListDto.addedBy[0]
          }}</v-avatar>
        </template>
      </v-tooltip>
      <div class="d-flex justify-end align-center mr-n2">
        <v-btn
          size="medium"
          icon="offline_bolt"
          variant="text"
          :disabled="props.songListDto.wasBoosted || props.songListDto.wasPlayed"
          @click="
            (event: Event) => props.onBoosted({ event: event, id: props.songListDto.id })
          "
          :color="isAdmin ? 'primary' : 'primary-container'"
        ></v-btn>
        <v-btn
          size="small"
          icon="keyboard_arrow_up"
          variant="text"
          :readonly="props.songListDto.wasDownVoted || props.songListDto.wasPlayed"
          :disabled="props.songListDto.wasUpVoted || props.songListDto.wasPlayed"
          @click="
            (event: Event) => props.onVotedUp({ event: event, id: props.songListDto.id })
          "
          color="primary"
        ></v-btn>
        <span class="outline">{{ props.songListDto.votes }}</span>
        <v-btn
          size="small"
          icon="keyboard_arrow_down"
          :readonly="props.songListDto.wasUpVoted || props.songListDto.wasPlayed"
          variant="text"
          :disabled="props.songListDto.wasDownVoted || props.songListDto.wasPlayed"
          @click="
            (event: Event) => props.onVotedDown({ event: event, id: props.songListDto.id })
          "
          color="primary"
        ></v-btn>
      </div>
    </div>
  </v-sheet>
</template>
<style scoped>
.outline {
  color: v-bind(outlineColor);
}
</style>
