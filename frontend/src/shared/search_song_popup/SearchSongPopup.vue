<script setup lang="ts">
import type { TSearchSongPopupEmits } from '@/shared/search_song_popup/TSearchSongPopupEmits.ts';
import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';
import MusicDataListElement from '@/shared/music_data_list_element/MusicDataListElement.vue';
import type IStringEvent from '../IStringEvent';
import { useDebounceFn } from '@vueuse/core';
import { SongRepository } from '@/infrastructure/songs/song_repository';
import type ISearchSongRequest from '@/infrastructure/songs/ISearchSongRequest';
import * as z from 'zod';
import type { IInfiniteScrollProps } from '../IInfiniteScrollProps';
import { shallowRef } from 'vue';

const popup = shallowRef(false);
const emit = defineEmits<TSearchSongPopupEmits>();
const songList = shallowRef<IMusicDataListElementDto[]>([]);
const nextPageToken = shallowRef<string | undefined>(undefined);
const isSearching = shallowRef(false);
const searchQueryString = shallowRef<string>('');
const hasNextPage = shallowRef(false);

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
    pageSize: 20,
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
    nextPageToken.value = response.meta?.nextPageToken;
  }
  isSearching.value = false;
}, 500);

function sendChosenSong(event: IStringEvent) {
  emit('on-song-choice', event.id);
}
async function loadNextPage({ done }: IInfiniteScrollProps) {
  if (nextPageToken.value) {
    const searchQuery: ISearchSongRequest = {
      query: searchQueryString.value,
      nextPageToken: nextPageToken.value,
      pageSize: 20,
    };
    const response = await SongRepository.searchSongs(searchQuery);
    if (response.isSuccess) {
      songList.value.push(
        ...response.data.map((song) => {
          return {
            id: song.videoId,
            author: song.author,
            title: song.title,
            imageUrl: song.albumCoverUrl,
          };
        }),
      );
      nextPageToken.value = response.meta?.nextPageToken;
      if (!response.meta?.hasNextPage) {
        hasNextPage.value = false;
        done('empty');
        return;
      }
    }
    hasNextPage.value = true;
    done('ok');
  }
  done('ok');
}
function resetSearch() {
  searchQueryString.value = '';
  songList.value = [];
  nextPageToken.value = undefined;
  hasNextPage.value = false;
}
</script>

<template>
  <v-dialog activator="parent" v-model="popup" @update:modelValue="resetSearch">
    <template v-slot:default>
      <v-row no-gutters justify="center">
        <v-col cols="12" sm="10" md="8" lg="6" xl="4">
          <v-card rounded="xl" color="surface">
            <v-text-field
              label="Search for song to add"
              prepend-inner-icon="search"
              append-inner-icon="close"
              hide-details
              max-length="50"
              min-length="2"
              :loading="isSearching"
              clearable
              v-model="searchQueryString"
              @update:modelValue="searchForSong"
              @click:append-inner="popup = false"
              single-line
              counter
            >
            </v-text-field>
            <v-infinite-scroll @load="loadNextPage" color="surface-container" height="60vh">
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
              <!-- Used to hide default scroll when nothing is being searched for -->
              <template v-if="!hasNextPage" v-slot:loading></template>
            </v-infinite-scroll>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-dialog>
</template>
