<script setup lang="ts">
import type { ISongListElementProps } from '@/room_page/ISongListElementProps.ts';
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import TouchscreenTooltip from '@/shared/touchscreen_tooltip/TouchscreenTooltip.vue';

const props = defineProps<ISongListElementProps>();
const backgroundColor = computed(() =>
  props.songListViewModel.wasPlayed ? 'outline-variant' : 'surface-container',
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
        :src="props.songListViewModel.albumCoverUrl"
      >
      </v-img>
    </div>
    <v-container
      class="w-0 flex-shrink-1 flex-grow-1 py-0 ma-0 mr-auto d-flex flex-column justify-between px-2"
    >
      <touchscreen-tooltip :open-on-hover="false"
                           :text="props.songListViewModel.title"
                           v-slot="{ tooltipProps }" >
        <span v-bind="tooltipProps" class="text-truncate">
            {{ props.songListViewModel.title }}</span
        >
      </touchscreen-tooltip>
      <touchscreen-tooltip :text="props.songListViewModel.author"
                           :open-on-hover="false"
                           v-slot="{ tooltipProps }">
        <span v-bind="tooltipProps" class="on-surface-variant text-truncate">{{
          props.songListViewModel.author
        }}</span>
      </touchscreen-tooltip>
    </v-container>
    <div class="d-flex flex-column justify-space-between align-end">
      <v-tooltip close-delay="2" open-on-click :text="props.songListViewModel.addedBy">
        <template v-slot:activator="{ props: tooltipProps }">
          <v-avatar v-bind="tooltipProps" size="x-small" color="primary">{{
            props.songListViewModel.addedBy[0]
          }}</v-avatar>
        </template>
      </v-tooltip>
      <div class="d-flex justify-end align-center mr-n2">
        <v-btn
          size="medium"
          icon="offline_bolt"
          variant="text"
          :disabled="props.songListViewModel.wasBoosted || props.songListViewModel.wasPlayed"
          @click="
            (event: Event) => props.onBoosted({ event: event, id: props.songListViewModel.id })
          "
          :color="isAdmin ? 'primary' : 'primary-container'"
        ></v-btn>
        <v-btn
          size="small"
          icon="keyboard_arrow_up"
          variant="text"
          :readonly="props.songListViewModel.wasDownVoted || props.songListViewModel.wasPlayed"
          :disabled="props.songListViewModel.wasUpVoted || props.songListViewModel.wasPlayed"
          @click="
            (event: Event) => props.onVotedUp({ event: event, id: props.songListViewModel.id })
          "
          color="primary"
        ></v-btn>
        <span class="outline">{{ props.songListViewModel.votes }}</span>
        <v-btn
          size="small"
          icon="keyboard_arrow_down"
          :readonly="props.songListViewModel.wasUpVoted || props.songListViewModel.wasPlayed"
          variant="text"
          :disabled="props.songListViewModel.wasDownVoted || props.songListViewModel.wasPlayed"
          @click="
            (event: Event) => props.onVotedDown({ event: event, id: props.songListViewModel.id })
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
