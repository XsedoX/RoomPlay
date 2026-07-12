import type { ISongListRoomStoreModel } from '@/infrastructure/room/ISongListRoomStoreModel.ts';
import type IGuidEvent from '@/shared/guid/IGuidEvent';

export interface ISongListElementProps {
  songListDto: ISongListRoomStoreModel;
  onVotedUp: (event: IGuidEvent) => void;
  onVotedDown: (event: IGuidEvent) => void;
  onBoosted: (event: IGuidEvent) => void;
  adminView: boolean;
}
