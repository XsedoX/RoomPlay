<script setup lang="ts">
import { shallowRef } from 'vue';
import MusicDataListElement from '@/shared/music_data_list_element/MusicDataListElement.vue';
import type ISearchSongPopupProps from '@/room_page/search_song_popup/ISearchSongPopupProps.ts';
import type { IMusicDataListElementViewModel } from '@/shared/music_data_list_element/IMusicDataListElementProps.ts';
import { Guid } from '@/utils/Guid.ts';

const dialog = shallowRef(false)
const props = defineProps<ISearchSongPopupProps>()
const songList: IMusicDataListElementViewModel[] = [
  {
    id: Guid.generate(),
    title: 'Bohemian Rhapsody',
    subtitle: 'Queen',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Stairway to Heaven',
    subtitle: 'Led Zeppelin',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Hotel California',
    subtitle: 'Eagles',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Bohemian Rhapsody',
    subtitle: 'Queen',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Stairway to Heaven',
    subtitle: 'Led Zeppelin',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Hotel California',
    subtitle: 'Eagles',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Bohemian Rhapsody',
    subtitle: 'Queen',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Stairway to Heaven',
    subtitle: 'Led Zeppelin',
    imageUrl: 'https://picsum.photos/200'
  },
  {
    id: Guid.generate(),
    title: 'Hotel California',
    subtitle: 'Eagles',
    imageUrl: 'https://picsum.photos/200'
  }
];
</script>

<template>
  <v-dialog activator="parent"
            v-model="dialog">
    <template v-slot:default>
      <v-row no-gutters justify="center">
        <v-col cols="12">
          <v-card rounded="xl" color="surface">
            <v-text-field
              label="Search for a Song"
              prepend-inner-icon="search"
              append-inner-icon="close"
              hide-details
              @click:append-inner="dialog = false"
              single-line>
            </v-text-field>
            <v-sheet color="surface-container" class="songs-list overflow-y-auto hide-scrollbar">
              <MusicDataListElement v-for="(song, index) in songList"
                                    :key="song.id.toString()"
                                    :class="{
                                      'rounded-b-xl': index === songList.length - 1,
                                    }"
                                    :musicDataListElementViewModel="song"
                                    :on-element-click="event=>props.chooseSong(event)">
              </MusicDataListElement>
            </v-sheet>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-dialog>
</template>
<style scoped>
@import '@/assets/shared.css';

.songs-list {
  max-height: 60vh;
}
</style>

