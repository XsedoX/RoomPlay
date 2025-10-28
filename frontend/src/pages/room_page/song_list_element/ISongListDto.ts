import type { IGuid } from '@/shared/Guid.ts';

export interface ISongListDto {
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
