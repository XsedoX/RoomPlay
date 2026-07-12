<script setup lang="ts">
import type { TSearchSongPopupEmits } from '@/shared/search_song_popup/TSearchSongPopupEmits.ts';
import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';
import { shallowRef } from 'vue';
import MusicDataListElement from '@/shared/music_data_list_element/MusicDataListElement.vue';
import type IStringEvent from '../IStringEvent';
import { useDebounceFn } from '@vueuse/core';
import { SongRepository } from '@/infrastructure/songs/song_repository';
import type ISearchSongRequest from '@/infrastructure/songs/ISearchSongRequest';
import * as z from 'zod';

const popup = shallowRef(false);
const emit = defineEmits<TSearchSongPopupEmits>();
const songList = shallowRef<IMusicDataListElementDto[]>([]);
const nextPageToken = shallowRef<string | undefined>(undefined);
const isSearching = shallowRef(false);

const querySchema = z
  .string()
  .min(2, 'Search query must be at least 2 characters long')
  .max(50, 'Search query must be at most 50 characters long');

const searchForSong = useDebounceFn(async (value: string) => {
  const validationResult = querySchema.safeParse(value);
  if (!validationResult.success) {
    songList.value = [];
    return;
  }
  const searchQuery: ISearchSongRequest = {
    query: value,
    nextPageToken: nextPageToken.value,
    pageSize: 10,
  };
  isSearching.value = true;
  const response = await SongRepository.searchSongs(searchQuery);
  if (response.isSuccess) {
    songList.value = [
      ...response.data.map((song) => {
        return {
          id: song.videoId,
          author: song.author,
          title: song.title,
          imageUrl: song.albumCoverUrl,
        };
      }),
    ];
    nextPageToken.value = response.meta!.nextPageToken;
  }
  isSearching.value = false;
}, 500);

function sendChosenSong(event: IStringEvent) {
  emit('on-song-choice', event.id);
}
</script>

<template>
  <v-dialog activator="parent" v-model="popup" @update:modelValue="songList = []">
    <template v-slot:default>
      <v-row no-gutters justify="center">
        <v-col cols="12" sm="10" md="8" lg="6" xl="4">
          <v-card rounded="xl" color="surface">
            <v-text-field
              label="Search for song to add"
              prepend-inner-icon="search"
              append-inner-icon="close"
              hide-details
              :loading="isSearching"
              clearable
              @update:modelValue="searchForSong"
              @click:append-inner="popup = false"
              single-line
              counter
            >
            </v-text-field>
            <v-sheet
              color="surface-container"
              class="music-list-popup-height overflow-y-auto hide-scrollbar"
            >
              <MusicDataListElement
                v-for="(song, index) in songList"
                :key="song.id"
                :class="{
                  'rounded-b-xl': index === songList.length - 1,
                }"
                :musicDataListDto="song"
                @on-music-choice="sendChosenSong"
              >
              </MusicDataListElement>
            </v-sheet>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-dialog>
</template>
