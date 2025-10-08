import type { ISongListViewModel } from '@/room_page/ISongListViewModel.ts';
import type IGuidEvent from '@/shared/IGuidEvent.ts';

export interface ISongListElementProps {
  songListViewModel: ISongListViewModel;
  onVotedUp: (event: IGuidEvent) => void;
  onVotedDown: (event: IGuidEvent) => void;
  onBoosted: (event: IGuidEvent) => void;
  adminView: boolean;
}
