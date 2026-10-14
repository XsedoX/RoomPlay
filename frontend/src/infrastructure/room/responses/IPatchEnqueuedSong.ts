import type { Guid } from '@/shared/guid/Guid';
import type { TSongState } from '../TSongState';
import type { TVoteStatus } from '../TVoteStatus';

export interface IPatchEnqueuedSong {
  title?: string;
  author?: string;
  addedBy?: string;
  votes?: number;
  albumCoverUrl?: string;
  id?: Guid;
  state?: TSongState;
  voteStatus?: TVoteStatus;
}
