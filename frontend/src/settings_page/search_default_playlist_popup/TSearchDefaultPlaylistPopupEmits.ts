import type IMusicDataListElementDto from '@/shared/music_data_list_element/IMusicDataListElementDto.ts';

export type TSearchDefaultPlaylistPopupEmits = {
  'on-playlist-choice': [payload: IMusicDataListElementDto];
}
