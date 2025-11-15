<script setup lang="ts">
import { computed, onMounted, shallowRef } from 'vue';
import SongListElement from '@/pages/room_page/song_list_element/SongListElement.vue';
import { type IGuid } from '@/shared/Guid.ts';
import SearchSongPopup from '@/shared/search_song_popup/SearchSongPopup.vue';
import TouchscreenTooltip from '@/shared/touchscreen_tooltip/TouchscreenTooltip.vue';
import type IGuidEvent from '@/shared/IGuidEvent.ts';
import SettingsMenu from '@/pages/room_page/settings_menu/SettingsMenu.vue';
import PageTitle from '@/shared/page_title/PageTitle.vue';
import { useUserStore } from '@/stores/user_store.ts';
import { useRoomStore } from '@/stores/room_store.ts';
import { useRouter } from 'vue-router';
import { Time } from '@/shared/Time.ts';
import { useIntervalFn } from '@vueuse/core';

const userStore = useUserStore();
const roomStore = useRoomStore();
const router = useRouter();
const playingSongMomentTimer = shallowRef<Time | null>(null);
const songTimerPercentage = computed(() => {
  if (!playingSongMomentTimer.value || !roomStore.playingSong) return 0;
  return Math.round(
    (playingSongMomentTimer.value!.totalSeconds() / roomStore.playingSong.lengthSeconds) * 100,
  );
});

onMounted(async () => {
  const isError = await roomStore.getRoom();
  if (isError) {
    await router.replace({ name: 'MainMenuPage' });
    return;
  }
  if (!roomStore.playingSong) {
    pause();
    return;
  }
  playingSongMomentTimer.value = Time.from(roomStore.playingSong.startedAtUtc).to();
  resume();
});
const { pause, resume } = useIntervalFn(() => {
  if (playingSongMomentTimer.value && roomStore.playingSong) {
    if (playingSongMomentTimer.value.totalSeconds() < roomStore.playingSong.lengthSeconds) {
      playingSongMomentTimer.value?.incrementSeconds();
    } else {
      pause();
    }
  }
}, 1000);
function onSongUpvoted(event: IGuidEvent) {
  roomStore.upVoteSong(event.id);
}
function onSongDownvoted(event: IGuidEvent) {
  roomStore.downVoteSong(event.id);
}
function onSongBoosted(event: IGuidEvent) {
  console.log('Song boosted:', event.id);
}
function chooseSong(id: IGuid) {
  console.log('Song chosen:', id);
}
async function leaveRoom() {
  await roomStore.leaveRoom();
  await router.replace({ name: 'MainMenuPage' });
}
</script>

<template>
  <v-container min-width="320px" fluid class="pa-0 pb-4 ma-0 h-screen d-flex flex-column">
    <v-row no-gutters>
      <v-col cols="12">
        <PageTitle
          :users-initial-letters="userStore.usersInitials"
          :title="roomStore.room?.name ?? 'I do not know what room this is...'"
          class="px-1"
        >
          <template v-slot:top-left-corner>
            <v-btn color="red" rounded="xl" @click="leaveRoom" variant="text"> Leave </v-btn>
          </template>
          <template v-slot:top-right-corner>
            <v-btn icon variant="text">
              <v-avatar size="small" color="primary">
                {{ userStore.usersInitials }}
              </v-avatar>
              <SettingsMenu></SettingsMenu>
            </v-btn>
          </template>
        </PageTitle>
      </v-col>
    </v-row>
    <v-row no-gutters justify="center">
      <v-col cols="11">
        <v-sheet
          v-if="roomStore.playingSong"
          color="surface-container"
          class="elevation-1 px-4 py-1 rounded-b-xl d-flex justify-start align-center ga-1"
        >
          <v-sheet color="transparent" class="d-flex w-0 flex-column flex-shrink-1 flex-grow-1">
            <v-sheet color="transparent" class="d-flex align-center overflow-hidden ga-1">
              <touchscreen-tooltip
                :text="roomStore.playingSong.title"
                :open-on-hover="false"
                v-slot="{ tooltipProps }"
              >
                <div
                  class="text-left text-subtitle-1 flex-grow-1 on-surface text-truncate"
                  v-bind="tooltipProps"
                >
                  {{ roomStore.playingSong.title }}
                </div>
              </touchscreen-tooltip>
              <touchscreen-tooltip
                :text="roomStore.playingSong.author"
                :open-on-hover="false"
                v-slot="{ tooltipProps }"
              >
                <div
                  class="text-body-1 flex-grow-1 on-surface-variant text-truncate text-right"
                  v-bind="tooltipProps"
                >
                  {{ roomStore.playingSong.author }}
                </div>
              </touchscreen-tooltip>
            </v-sheet>
            <v-sheet color="transparent" class="d-flex align-center justify-start ga-2">
              <div class="on-surface-variant">
                {{ playingSongMomentTimer!.toString() }}
              </div>
              <v-progress-linear
                height="6"
                color="primary"
                :model-value="songTimerPercentage"
              ></v-progress-linear>
              <div class="on-surface-variant">
                {{ Time.fromSeconds(roomStore.playingSong.lengthSeconds).toString() }}
              </div>
            </v-sheet>
          </v-sheet>
          <v-btn
            v-if="roomStore.isHost"
            icon="skip_next"
            variant="text"
            class="mr-n3"
            color="black"
          ></v-btn>
        </v-sheet>
      </v-col>
    </v-row>
    <v-row class="py-2 overflow-y-hidden h-100" no-gutters justify="center">
      <v-col cols="11" class="fill-height">
        <v-sheet class="fill-height overflow-y-auto hide-scrollbar" rounded="xl">
          <div v-if="roomStore.songs != null">
            <div v-for="(song, index) in roomStore.songs" :key="song.id.toString()">
              <song-list-element
                :class="{
                  'rounded-t-xl': index === 0,
                  'rounded-b-xl': index === roomStore.songs.length - 1,
                }"
                :onVotedDown="onSongDownvoted"
                :onBoosted="onSongBoosted"
                :onVotedUp="onSongUpvoted"
                :songListDto="song"
                :adminView="roomStore.isHost"
              >
              </song-list-element>
              <v-divider v-if="index < roomStore.songs.length - 1"></v-divider>
            </div>
          </div>
        </v-sheet>
      </v-col>
    </v-row>
    <v-row no-gutters justify="center">
      <v-col cols="11">
        <v-text-field
          label="Search"
          prepend-inner-icon="search"
          hide-details
          rounded="xl"
          variant="solo-filled"
          single-line
        >
          <search-song-popup @on-song-choice="chooseSong"></search-song-popup>
        </v-text-field>
      </v-col>
    </v-row>
  </v-container>
</template>
<style scoped>
@import '@/assets/shared.css';
</style>
