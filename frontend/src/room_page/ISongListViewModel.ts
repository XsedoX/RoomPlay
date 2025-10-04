import type { IGuid } from '@/utils/Guid.ts';

export interface ISongListViewModel {
  title: string;
  author: string;
  addedBy: string;
  votes: number;
  albumCoverUrl: string;
  id: IGuid;
  wasPlayed: boolean;
  wasUpVoted: boolean;
  wasDownVoted: boolean;
  wasBoosted: boolean;
}
