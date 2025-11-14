
export interface IGetRoomResponse {
  name: string;
  boostData: IGetRoomBoostData | null;
  qrCode: string;
  userRole: string;
  songs: IGetRoomSongListResponse[];
  playingSong: IGetRoomPlayingSongResponse | null;
}
export interface IGetRoomSongListResponse {
  title: string;
  author: string;
  addedBy: string;
  votes: number;
  albumCoverUrl: string;
  id: string;
  state: string;
  voteStatus: string;
}
export interface IGetRoomPlayingSongResponse {
  title: string;
  author: string;
  startedAtUtc: Date;
  lengthSeconds: number;
}
export interface IGetRoomBoostData {
  boostUsedAtUtc: Date;
  boostCooldownSeconds: number;
}
