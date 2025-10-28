import type { ISongListDto } from '@/pages/room_page/song_list_element/ISongListDto.ts';
import type IGuidEvent from '@/shared/IGuidEvent.ts';

export interface ISongListElementProps {
  songListDto: ISongListDto;
  onVotedUp: (event: IGuidEvent) => void;
  onVotedDown: (event: IGuidEvent) => void;
  onBoosted: (event: IGuidEvent) => void;
  adminView: boolean;
}
