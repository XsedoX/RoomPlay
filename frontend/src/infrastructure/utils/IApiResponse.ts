import type { TIncomingActionName } from '../websockets/TIncomingActionName';

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

export interface IApiWebSocketResponse<T = void> {
  payload: IApiSuccessResponse<T> | IApiProblemDetailsResponse;
  actionName: TIncomingActionName;
}

interface IApiJsonPatchOperation<T = void> {
  op: 'add' | 'remove' | 'replace' | 'move' | 'copy' | 'test';
  path: string;
  value: T;
}

export type IApiJsonPatchResponse<T = void> = IApiSuccessResponse<IApiJsonPatchOperation<T>[]>;
