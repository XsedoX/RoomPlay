export default interface ISearchSongRequest {
  query: string;
  nextPageToken?: string | undefined;
  pageSize: number;
}
