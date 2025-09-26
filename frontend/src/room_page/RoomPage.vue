<script setup lang="ts">
import { type Ref, ref } from 'vue';
import SongListElement from '@/room_page/SongListElement.vue';
import type { ISongListViewModel } from '@/room_page/ISongListViewModel.ts';
import type { SongListEvent } from '@/room_page/ISongListElementProps.ts';
import { Guid } from '@/utils/Guid.ts';

const isAdmin = ref(true);
const exampleSongs: Ref<ISongListViewModel[]> = ref([
  {
    title: 'Bohemian Rhapsody',
    author: 'Queen',
    addedBy: 'Alice',
    votes: 5,
    albumCoverUrl: 'https://picsum.photos/200',
    id: Guid.generate(),
    wasPlayed: true,
    wasBoosted: false,
    wasUpVoted: false,
    wasDownVoted: false
  },
  {
    title: 'Imagine',
    author: 'John Lennon',
    addedBy: 'Bob',
    votes: 3,
    albumCoverUrl: 'https://picsum.photos/200',
    id: Guid.generate(),
    wasPlayed: false,
    wasBoosted: false,
    wasUpVoted: false,
    wasDownVoted: true
  },
  {
    title: 'Billie Jean',
    author: 'Michael Jackson',
    addedBy: 'Charlie',
    votes: 7,
    albumCoverUrl: 'https://picsum.photos/200',
    id: Guid.generate(),
    wasPlayed: false,
    wasBoosted: false,
    wasUpVoted: true,
    wasDownVoted: false
  },
]);

function onSongUpvoted(event: SongListEvent) {
  const song = exampleSongs.value.find((song) => song.id === event.id);
  if (song) {
    song.votes += 1;
  }
}
function onSongDownvoted(event: SongListEvent) {
  const song = exampleSongs.value.find((song) => song.id === event.id);
  if (song && song.votes > 0) {
    song.votes -= 1;
  }
}
function onSongBoosted(event: SongListEvent) {
  console.log('Song boosted:', event.id);
}
</script>

<template>
  <v-container
    min-width="320px"
    fluid
    class="pa-0 ma-0 fill-height align-content-start"
    style="display: grid; grid-template-rows: auto auto 1fr auto"
  >
    <v-row no-gutters>
      <v-col cols="12">
        <v-sheet
          color="surface-container-highest"
          class="elevation-4 pa-1 d-flex align-center justify-space-between"
        >
          <div class="w-25 d-flex justify-start">
            <v-btn color="red" rounded="xl" variant="text"> Leave </v-btn>
          </div>
          <span class="text-h6 text-no-wrap">Room Name</span>
          <div class="w-25 d-flex justify-end">
            <v-btn icon="qr_code" variant="text"></v-btn>
            <v-btn v-if="isAdmin" icon="settings" variant="text"></v-btn>
          </div>
        </v-sheet>
      </v-col>
    </v-row>
    <v-row no-gutters justify="center">
      <v-col cols="11">
        <v-sheet color="surface-container" class="elevation-1 px-4 py-1 rounded-b-xl">
          <v-row no-gutters>
            <v-col>
              <v-row no-gutters justify="space-between" class="align-center">
                <v-col cols="auto">
                  <span class="text-subtitle-1 on-surface">Song Title</span>
                </v-col>
                <v-col cols="auto">
                  <span class="text-body-1 on-surface-variant">Author</span>
                </v-col>
              </v-row>
              <v-row no-gutters class="align-center ga-2">
                <v-col cols="auto">
                  <span class="on-surface-variant">0:45</span>
                </v-col>
                <v-col class="d-flex align-center">
                  <v-progress-linear
                    height="6"
                    color="primary"
                    model-value="20"
                    class="w-100"
                  ></v-progress-linear>
                </v-col>
                <v-col cols="auto">
                  <span class="on-surface-variant">3:45</span>
                </v-col>
              </v-row>
            </v-col>
            <v-col v-if="isAdmin" cols="auto" class="align-center d-flex mr-n3">
              <v-btn icon="skip_next" variant="text" color="black"></v-btn>
            </v-col>
          </v-row>
        </v-sheet>
      </v-col>
    </v-row>
    <v-row class="pt-2 align-self-start" no-gutters justify="center">
      <v-col cols="11">
        <v-sheet class="fill-height overflow-hidden song-list" rounded="xl">
          <v-sheet v-for="song in exampleSongs" :key="song.id.toString()">
            <song-list-element
              :onVotedDown="onSongDownvoted"
              :onBoosted="onSongBoosted"
              :onVotedUp="onSongUpvoted"
              :songListViewModel="song"
              :adminView="isAdmin">
            </song-list-element>
          </v-sheet>
        </v-sheet>
      </v-col>
    </v-row>
    <v-row>
      <v-col> </v-col>
    </v-row>
  </v-container>
</template>
<style scoped>
  .song-list{
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
  }
</style>
