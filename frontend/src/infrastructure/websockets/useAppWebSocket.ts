import { useWebSocket } from '@vueuse/core';
import type { IApiWebSocketResponse } from '../utils/IApiResponse';
import { apiWebSocketDispatcher } from './ApiWebSocketActionDispatcher';
import { useNotificationStore } from '@/stores/notification_store';
import { TSnackbarColor } from '../utils/TSnackbarColor';
import type { TOutgoingActionName } from './TOutgoingActionName';
import { shallowRef } from 'vue';
import { AuthenticationRepository } from '../authentication/authentication_repository';
import { useUserStore } from '@/stores/user_store';
import router from '@/router';

type RefreshResult = Awaited<ReturnType<typeof AuthenticationRepository.refreshToken>>;

// Module-scoped so several live sockets never refresh at the same time.
const inFlight = shallowRef<Promise<RefreshResult | null> | null>(null);
/**
 * Refreshes the JWT at most once while a refresh is already running.
 * Returns `null` on a transport error (not an auth failure).
 */
function refreshTokenOnce(): void {
  if (!inFlight.value) {
    inFlight.value = AuthenticationRepository.refreshToken()
      .catch(() => null)
      .finally(() => {
        inFlight.value = null;
      });
  }
}

async function forceLogout() {
  const userStore = useUserStore();
  await userStore.logout();
  await router.replace({ name: 'LoginPage' });
}

export function useAppWebSocket(upgradeToWebSocketEndpoint: string) {
  const url = import.meta.env['VITE_API_BASE_URL'] + upgradeToWebSocketEndpoint;
  // At most one refresh per disconnect round, so an outage can't hammer /refresh.
  let refreshAttemptedThisRound = false;

  const {
    send: notWrappedSend,
    close,
    status,
  } = useWebSocket(url, {
    autoReconnect: true,
    autoConnect: true,
    onConnected: (ws) => {
      console.log('WebSocket connected:', ws);
      refreshAttemptedThisRound = false;
    },
    onMessage: (_, event) => {
      try {
        const message = JSON.parse(event.data) as IApiWebSocketResponse<never>;
        console.log('WebSocket message received:', message);

        apiWebSocketDispatcher.dispatch(message.actionName, message.payload);
      } catch (error) {
        const notificationStore = useNotificationStore();
        notificationStore.showSnackbar(
          'Something went really badly try refreshing the page or reloging.',
          TSnackbarColor.ERROR,
        );
      }
    },
    onDisconnected: (ws, event) => {
      console.log('WebSocket disconnected:', ws, event);
    },
    onError: async (_, event) => {
      console.log(
        'WebSocket error occurred. Attempting to refresh token if not already done this round.',
        event,
      );
      if (!refreshAttemptedThisRound) {
        refreshAttemptedThisRound = true;
        refreshTokenOnce();
        const response = await inFlight.value;

        if (response && !response.isSuccess) {
          // Definitive auth failure (refresh token gone/expired) -> bail out.
          console.log('WebSocket disconnected due to auth failure. Forcing logout.');
          await forceLogout();
          return;
        }
        // response === null: transport error, fall through and retry later.
      }
    },
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
  return {
    send,
    close,
  };
}
