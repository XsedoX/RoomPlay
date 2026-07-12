import type { TSongState } from '@/infrastructure/room/TSongState.ts';
import type { TVoteStatus } from '@/infrastructure/room/TVoteStatus.ts';
import type { IGuid } from '@/shared/guid/Guid';

export interface ISongListRoomStoreModel {
  title: string;
  author: string;
  addedBy: string;
  votes: number;
  albumCoverUrl: string;
  id: IGuid;
  state: TSongState;
  voteStatus: TVoteStatus;
}
