import type { IPageMetaDto } from './IApiResponse';

interface IRepositorySuccessResponse<T> {
  data: T;
  meta?: IPageMetaDto | undefined;
  isSuccess: true;
}
interface IRepositoryFailureResponse {
  isSuccess: false;
  validationErrors?: Record<string, string> | undefined;
}

export type IRepositoryResponse<T = void> =
  IRepositorySuccessResponse<T> | IRepositoryFailureResponse;
