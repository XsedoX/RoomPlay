<script setup lang="ts">
import { type Ref, ref } from 'vue';
import SongListElement from '@/room_page/SongListElement.vue';
import type { ISongListViewModel } from '@/room_page/ISongListViewModel.ts';
import type { SongListEvent } from '@/room_page/ISongListElementProps.ts';
import { Guid } from '@/utils/Guid.ts';
import SearchSongPopup from '@/room_page/SearchSongPopup.vue';

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
  },{
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
  },{
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
  },{
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
    class="pa-0 pb-4 ma-0 h-screen d-flex flex-column">
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
        <v-sheet color="surface-container"
                 class="elevation-1 px-4 py-1 rounded-b-xl d-flex justify-start align-center ga-1">
          <v-sheet color="transparent"
                   class="d-flex w-0 flex-column flex-shrink-1 flex-grow-1">
            <v-sheet color="transparent"
                     class="d-flex align-center overflow-hidden ga-1">
              <div class="text-left text-subtitle-1 flex-grow-1 on-surface text-truncate">
                Song Titledwdwdwdwdwdwdwdwd
              </div>
              <div class="text-body-1 flex-grow-1 on-surface-variant text-truncate text-right">
                Authordawwdawdawdwadwa
              </div>
            </v-sheet>
            <v-sheet color="transparent"
                     class="d-flex align-center justify-start ga-2">
              <div class="on-surface-variant">
                0:45
              </div>
              <v-progress-linear
                height="6"
                color="primary"
                model-value="20"
              ></v-progress-linear>
              <div class="on-surface-variant">
                3:45
              </div>
            </v-sheet>
          </v-sheet>
          <v-btn v-if="isAdmin" icon="skip_next" variant="text" class="mr-n3" color="black"></v-btn>
        </v-sheet>
      </v-col>
    </v-row>
    <v-row class="py-2 overflow-y-hidden" no-gutters justify="center">
      <v-col cols="11" class="fill-height">
        <v-sheet class="fill-height overflow-y-auto hide-scrollbar" rounded="xl">
          <song-list-element
            v-for="(song, index) in exampleSongs"
            :class="{
              'rounded-t-xl': index === 0,
              'rounded-b-xl': index === exampleSongs.length - 1,
            }"
            :key="song.id.toString()"
            :onVotedDown="onSongDownvoted"
            :onBoosted="onSongBoosted"
            :onVotedUp="onSongUpvoted"
            :songListViewModel="song"
            :adminView="isAdmin">
          </song-list-element>
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
          color="on-surface-variant"
        >
          <SearchSongPopup></SearchSongPopup>
        </v-text-field>
      </v-col>
    </v-row>
  </v-container>
</template>
<style scoped>
.hide-scrollbar::-webkit-scrollbar {
  display: none;
}

.hide-scrollbar {
  -ms-overflow-style: none; /* IE and Edge */
  scrollbar-width: none; /* Firefox */
}
</style>
