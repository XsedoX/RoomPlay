import type { IGuid } from '@/shared/Guid.ts';
import type { TSongState } from '@/infrastructure/room/TSongState.ts';
import type { TVoteStatus } from '@/infrastructure/room/TVoteStatus.ts';

export interface ISongListRoomStoreModel {
  title: string;
  author: string;
  addedBy: string;
  votes: number;
  albumCoverUrl: string;
  id: IGuid;
  state: TSongState
  voteStatus: TVoteStatus
}
