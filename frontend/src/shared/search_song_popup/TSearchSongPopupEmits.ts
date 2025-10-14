import { type IGuid } from '@/shared/Guid.ts';

export type TSearchSongPopupEmits = {
  'on-song-choice': [id: IGuid];
}
