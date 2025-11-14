import type { IGuid } from '@/shared/Guid.ts';
import type { TSongState } from '@/infrastructure/room/TSongState.ts';

export interface ISongListRoomStoreModel {
  title: string;
  author: string;
  addedBy: string;
  votes: number;
  albumCoverUrl: string;
  id: IGuid;
  state: TSongState
}
