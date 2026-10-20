import type { IApiProblemDetailsResponse } from './IApiResponse';

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number | undefined, // undefined = request never left / timed out
    public readonly problemDetails: IApiProblemDetailsResponse | undefined,
    cause?: unknown,
  ) {
    super(message);
    this.cause = cause;
  }
}
