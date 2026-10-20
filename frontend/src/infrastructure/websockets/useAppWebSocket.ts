import { useWebSocket } from '@vueuse/core';
import type { IApiSuccessResponse, IApiWebSocketResponse } from '../utils/IApiResponse';
import { apiWebSocketDispatcher } from './ApiWebSocketActionDispatcher';
import type { TOutgoingActionName } from './TOutgoingActionName';
import { onUnmounted, shallowRef } from 'vue';
import api_client from '../utils/api_client';
import type { ApiError } from '../utils/ApiError';

interface IGetWebSocketTokenResponse {
  token: string;
}

export function useAppWebSocket(
  upgradeToWebSocketEndpoint: string,
  getWebSocketTicketEndpoint: string,
) {
  const upgradeUrl = shallowRef('');
  const isClosedPermanently = shallowRef(false);
  const reconnectAttempt = shallowRef(0);
  const reconnectDelay = shallowRef(1);

  async function fetchTicketAndSetUrl() {
    if (isClosedPermanently.value) return;
    await api_client
      .get<IApiSuccessResponse<IGetWebSocketTokenResponse>>(getWebSocketTicketEndpoint, {
        silent: true,
      })
      .then((response) => {
        if (isClosedPermanently.value) return;
        upgradeUrl.value = `${import.meta.env['VITE_API_WS_URL']}${upgradeToWebSocketEndpoint}?token=${response.data.data.token}`;
        open();
      })
      .catch((error: ApiError) => {
        console.error('Error fetching WebSocket upgrade ticket:', error);
        upgradeUrl.value = '';
        if (error.problemDetails?.status !== 401) {
          scheduleReconnect();
        }
      });
  }

  let reconnectTimer: ReturnType<typeof setTimeout> | undefined;

  function clearReconnectTimer() {
    if (reconnectTimer !== undefined) {
      clearTimeout(reconnectTimer);
      reconnectTimer = undefined;
    }
  }

  function scheduleReconnect() {
    if (isClosedPermanently.value) return;
    // Ensure only one pending reconnect at a time.
    if (reconnectTimer !== undefined) return;
    reconnectDelay.value = Math.min(1000 * 2 ** reconnectAttempt.value, 30000);
    reconnectAttempt.value++;
    reconnectTimer = setTimeout(() => {
      reconnectTimer = undefined;
      fetchTicketAndSetUrl();
    }, reconnectDelay.value);
  }
  const {
    send: notWrappedSend,
    close,
    open,
    status,
    ws: currentWs,
  } = useWebSocket(upgradeUrl, {
    // Don't let VueUse auto-open on URL change; we open manually after fetching a ticket.
    autoConnect: false,
    immediate: false,
    onConnected: (ws) => {
      reconnectAttempt.value = 0; // Reset reconnect attempts on successful connection
      console.log('WebSocket connected:', ws);
    },
    onMessage: (_, event) => {
      try {
        const message = JSON.parse(event.data) as IApiWebSocketResponse<never>;
        console.log('WebSocket message received:', message);

        apiWebSocketDispatcher.dispatch(message.actionName, message.payload);
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    },
    onDisconnected: (ws, event) => {
      console.log('WebSocket disconnected:', ws, event);
      // Ignore close events from sockets that were already replaced by a newer one.
      if (currentWs.value && currentWs.value !== ws) return;
      if (!isClosedPermanently.value) {
        scheduleReconnect();
      }
    },
  });
  fetchTicketAndSetUrl();
  onUnmounted(() => {
    closePermanently();
  });
  function send<TPayload>(actionName: TOutgoingActionName, payload: TPayload) {
    if (status.value === 'OPEN') {
      const message = {
        actionName: actionName,
        payload: payload,
      };
      console.log(JSON.stringify(message));
      notWrappedSend(JSON.stringify(message));
    } else {
      console.error('WebSocket is not open. Cannot send message.');
    }
  }
  function closePermanently() {
    upgradeUrl.value = '';
    isClosedPermanently.value = true;
    clearReconnectTimer();
    close();
  }
  return {
    send,
    close: closePermanently,
    status,
    reconnectAttempt,
    reconnectDelay,
  };
}
