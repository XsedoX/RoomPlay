<script setup lang="ts">
import { Guid } from '@/shared/Guid.ts';
import SearchMusicPopupBase from '@/shared/search_popup_base/SearchMusicPopupBase.vue';
import type IGuidEvent from '@/shared/IGuidEvent.ts';
import type { TSearchSongPopupEmits } from '@/room_page/search_song_popup/TSearchSongPopupEmits.ts';
import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';

const emit = defineEmits<TSearchSongPopupEmits>();
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
    id: Guid.generate(),
    title: 'Hotel California',
    subtitle: 'Eagles',
    imageUrl: 'https://picsum.photos/200',
  },
];
function sendChosenSong(event: IGuidEvent) {
  emit('on-song-choice', event.id);
}
</script>

<template>
  <search-music-popup-base
    search-box-placeholder="Search for a Song"
    :music-list="songList"
    @on-music-choice="sendChosenSong"
  ></search-music-popup-base>
</template>
