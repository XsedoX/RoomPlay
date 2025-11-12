
export interface IGetRoomResponse {
  name: string;
  boostUsedAtUtc: Date | null;
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
}
export interface IGetRoomPlayingSongResponse {
  title: string;
  author: string;
  startedAtUtc: Date;
  lengthSeconds: number;
}
