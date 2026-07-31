export interface IWebSocketEnvelope<T> {
  ActionName: string;
  Payload: T;
}
