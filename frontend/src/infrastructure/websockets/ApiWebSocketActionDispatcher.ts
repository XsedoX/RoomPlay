import type { TIncomingActionName } from './TIncomingActionName';

type ApiWebSocketHandler<TResponse> = (response: TResponse) => void;
type InternalHandler = (response: unknown) => void;

export class ApiWebSocketActionDispatcher {
  private handlers: Map<TIncomingActionName, Set<InternalHandler>> = new Map();

  public on<TResponse>(
    actionName: TIncomingActionName,
    handler: ApiWebSocketHandler<TResponse>,
  ): () => void {
    if (!this.handlers.has(actionName)) {
      this.handlers.set(actionName, new Set());
    }

    const handlerSet = this.handlers.get(actionName)!;
    const internalHandler = handler as InternalHandler;
    handlerSet.add(internalHandler);

    // Return an unsubscription callback
    return () => {
      handlerSet.delete(internalHandler);
      if (handlerSet.size === 0) {
        this.handlers.delete(actionName);
      }
    };
  }

  public dispatch<TResponse>(actionName: TIncomingActionName, response: TResponse): void {
    const handlerSet = this.handlers.get(actionName);
    if (handlerSet) {
      handlerSet.forEach((handler) => handler(response));
    }
  }
}

export const apiWebSocketDispatcher = new ApiWebSocketActionDispatcher();
