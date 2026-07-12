export interface IPageMetaDto {
  nextPageToken?: string;
  prevPageToken?: string;
  pageSize: number;
  hasNextPage: boolean;
}
export interface IApiSuccessResponse<T = void> {
  meta?: IPageMetaDto;
  data: T;
}
export interface IApiProblemDetailsResponse {
  type: string;
  title: string;
  description: string;
  instance: string;
  status: number;
  validationErrors?: Record<string, string>;
}
