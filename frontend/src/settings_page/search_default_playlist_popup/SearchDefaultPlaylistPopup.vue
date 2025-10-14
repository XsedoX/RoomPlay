<script setup lang="ts">
import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';
import SearchMusicPopupBase from '@/shared/search_popup_base/SearchMusicPopupBase.vue';
import type { TSearchDefaultPlaylistPopupEmits } from '@/settings_page/search_default_playlist_popup/TSearchDefaultPlaylistPopupEmits.ts';
import type IGuidEvent from '@/shared/IGuidEvent.ts';
import { Guid } from '@/shared/Guid.ts';

const emit = defineEmits<TSearchDefaultPlaylistPopupEmits>();
const playlistList: IMusicDataListElementDto[] = [
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 1',
    subtitle: '15 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 2',
    subtitle: '8 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 3',
    subtitle: '22 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 1',
    subtitle: '15 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 2',
    subtitle: '8 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 3',
    subtitle: '22 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 1',
    subtitle: '15 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 2',
    subtitle: '8 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 3',
    subtitle: '22 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 1',
    subtitle: '15 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 2',
    subtitle: '8 songs',
    imageUrl: 'https://picsum.photos/200',
  },
  {
    id: Guid.generate(),
    title: 'Sample Playlist Title 3',
    subtitle: '22 songs',
    imageUrl: 'https://picsum.photos/200',
  },
];
function sendChosenPlaylist(event: IGuidEvent) {
  const chosenPlaylist = playlistList.find((playlist) => playlist.id === event.id);
  if(!chosenPlaylist) throw new Error(
    `Could not find playlist with id ${event.id} in playlistList: ${playlistList}`
  )
  emit('on-playlist-choice', chosenPlaylist);
}
</script>

<template>
  <search-music-popup-base
    search-box-placeholder="Search for a Playlist"
    :music-list="playlistList"
    @on-music-choice="sendChosenPlaylist"
  ></search-music-popup-base>
</template>
