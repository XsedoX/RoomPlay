<script setup lang="ts">
import SettingsListElement from '@/pages/settings_page/settings_list_element/SettingsListElement.vue';
import { ref, shallowRef } from 'vue';
import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';
import SearchSongPopup from '@/shared/search_song_popup/SearchSongPopup.vue';
import { Guid, type IGuid } from '@/shared/Guid.ts';
import { type ITime, Time } from '@/shared/Time.ts';

const songList: IMusicDataListElementDto[] = [
  {
    id: Guid.generate(),
    title: 'Bohemian Rhapsody',
    subtitle: 'Queen',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Stairway to Heaven',
    subtitle: 'Led Zeppelin',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Hotel California',
    subtitle: 'Eagles',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Bohemian Rhapsody',
    subtitle: 'Queen',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Stairway to Heaven',
    subtitle: 'Led Zeppelin',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Hotel Californiadwwwwwwwwwwwwwwwwwwwwwwwww',
    subtitle: 'Eaglesdwwwwwwwwwwwwwwwwwwwwwwwwwwwwww',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Bohemian Rhapsody',
    subtitle: 'Queen',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Stairway to Heaven',
    subtitle: 'Led Zeppelin',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: new Guid('9593de37-416c-4bc8-9eee-8a98b7db66c7'),
    title: 'Hotel California',
    subtitle: 'Eagles',
    imageUrl: 'https://picsum.photos/200',
  },
];
const chosenSong = ref<IMusicDataListElementDto | null>(null);
function chooseSong(id: IGuid) {
  chosenSong.value = songList.find((song) => song.id.toString() === id.toString())!;
}
const rapidSongPlayTime = ref<ITime>(new Time());
const showTimePicker = shallowRef(false);
</script>

<template>
  <v-list bg-color="transparent" density="compact" width="100%">
    <settings-list-element
      header="CHOOSE A SONG"
      :hint="chosenSong?.subtitle"
      :sub-header="chosenSong?.title ?? 'Choose a song.'"
    >
      <template v-slot:image>
        <v-img
          cover
          v-if="chosenSong"
          max-width="50px"
          min-width="50px"
          rounded="lg"
          aspect-ratio="1/1"
          :src="chosenSong.imageUrl"
        ></v-img>
      </template>
      <template v-slot:default>
        <v-btn variant="text" color="primary" rounded="xl">
          Change
          <search-song-popup @on-song-choice="chooseSong"></search-song-popup>
        </v-btn>
      </template>
    </settings-list-element>
    <settings-list-element v-if="chosenSong" header="SET TIME" sub-header="Time to play:">
      <v-btn color="primary" variant="text" append-icon="arrow_drop_down">
        {{ rapidSongPlayTime }}
        <v-dialog v-model="showTimePicker" activator="parent" width="auto">
          <v-time-picker
            format="24hr"
            color="primary"
            :model-value="rapidSongPlayTime.toString()"
            @update:model-value="(value)=>rapidSongPlayTime = new Time(value)"
          ></v-time-picker>
        </v-dialog>
      </v-btn>
    </settings-list-element>
  </v-list>
</template>
<style scoped>
:deep(.v-list-item--density-compact:not(.v-list-item--nav).v-list-item--one-line) {
  padding-right: 0;
  padding-left: 0;
}
:deep(.remove-padding) {
  padding-inline-start: 8px !important;
}
</style>
