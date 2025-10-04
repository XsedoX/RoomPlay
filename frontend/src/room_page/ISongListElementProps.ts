import type { ISongListViewModel } from '@/room_page/ISongListViewModel.ts';
import type { IGuid } from '@/utils/Guid.ts';

export interface ISongListElementProps {
  songListViewModel: ISongListViewModel;
  onVotedUp: (event: SongListEvent) => void;
  onVotedDown: (event: SongListEvent) => void;
  onBoosted: (event: SongListEvent) => void;
  adminView: boolean;
}
export interface SongListEvent{
  event: Event;
  id: IGuid;
}
