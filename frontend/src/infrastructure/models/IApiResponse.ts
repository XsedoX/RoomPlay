export default interface IApiResponse<T> {
  data: T|undefined;
  message: string;
}
