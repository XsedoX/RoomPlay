interface IApiSuccessResponse<T> {
  data: T;
  isSuccess: true;
}
interface IApiFailureResponse {
  isSuccess: false;
  type: string;
  title: string;
  description: string;
  instance: string;
  status: number;
  validationErrors?: Record<string, string>;
}
export type IApiResponse<T = void> = IApiSuccessResponse<T> | IApiFailureResponse;
