<script setup lang="ts">
import type { ISongListElementProps } from '@/pages/room_page/song_list_element/ISongListElementProps.ts';
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import TouchscreenTooltip from '@/shared/touchscreen_tooltip/TouchscreenTooltip.vue';
import { useRoomStore } from '@/stores/room_store.ts';
import { TSongState } from '@/infrastructure/room/TSongState.ts';
import { TVoteStatus } from '@/infrastructure/room/TVoteStatus.ts';

const props = defineProps<ISongListElementProps>();
const roomStore = useRoomStore();
const backgroundColor = computed(() => (wasPlayed.value ? 'outline-variant' : 'surface-container'));
const wasPlayed = computed(() => props.songListDto.state === TSongState.played);
const wasUpVoted = computed(() => props.songListDto.voteStatus === TVoteStatus.upvoted);
const wasDownVoted = computed(() => props.songListDto.voteStatus === TVoteStatus.downvoted);
const theme = useTheme();
const outlineColor = theme.current.value.colors['outline'];
const isAdmin = true;
</script>

<template>
  <v-sheet
    data-testid="song-list-element"
    class="d-flex px-4 py-2 justify-start align-center w-100"
    :color="backgroundColor"
  >
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
      <touchscreen-tooltip
        :open-on-hover="false"
        :text="props.songListDto.title"
        v-slot="{ tooltipProps }"
      >
        <span v-bind="tooltipProps" class="text-truncate"> {{ props.songListDto.title }}</span>
      </touchscreen-tooltip>
      <touchscreen-tooltip
        :text="props.songListDto.author"
        :open-on-hover="false"
        v-slot="{ tooltipProps }"
      >
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
          :disabled="!roomStore.isBoostAvailable || wasPlayed"
          @click="(event: Event) => props.onBoosted({ event: event, id: props.songListDto.id })"
          :color="isAdmin ? 'primary' : 'primary-container'"
        ></v-btn>
        <v-btn
          size="small"
          icon="keyboard_arrow_up"
          variant="text"
          :disabled="wasPlayed || wasDownVoted"
          :readonly="wasUpVoted || wasPlayed || wasDownVoted"
          @click="(event: Event) => props.onVotedUp({ event: event, id: props.songListDto.id })"
          color="primary"
        ></v-btn>
        <span class="outline">{{ props.songListDto.votes }}</span>
        <v-btn
          size="small"
          icon="keyboard_arrow_down"
          :disabled="wasPlayed || wasUpVoted"
          variant="text"
          :readonly="wasDownVoted || wasPlayed || wasUpVoted"
          @click="(event: Event) => props.onVotedDown({ event: event, id: props.songListDto.id })"
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
