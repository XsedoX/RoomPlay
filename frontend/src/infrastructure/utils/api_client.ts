import axios, { type AxiosInstance } from 'axios';
import { AuthenticationService } from '@/infrastructure/authentication/authentication_service.ts'
import { useUserStore } from '@/stores/user_store.ts';
import { PlatformDiscoverer } from '@/infrastructure/utils/platform_discoverer.ts';
import router from '@/router';



const axiosParams = {
  baseURL: import.meta.env["VITE_API_BASE_URL"],
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    'Origin': import.meta.env["VITE_APP_ORIGIN"],
    'X-Device-Type': PlatformDiscoverer.getDeviceType()
  },
  timeout: 10000,
  withCredentials: true,
};


const api_client = axios.create(axiosParams)

const api = (axiosInstance: AxiosInstance) => {
  return {
    get: <T>(url: string, config = {}) => axiosInstance.get<T>(url, config),
    post: <T>(url: string, body?: unknown, config = {}) => axiosInstance.post<T>(url, body, config),
    delete: <T>(url: string, config = {}) => axiosInstance.delete<T>(url, config),
  };
};

function createRefreshTokenInterceptor(){
  const interceptor = api_client.interceptors.response.use(
    response => response,
    async error => {
      const originalRequest = error.config;
      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error);
      }
      originalRequest._retry = true;
      api_client.interceptors.response.eject(interceptor);

      try {
        const success = await AuthenticationService.refreshToken();
        if (!success) {
          const userStore = useUserStore();
          await userStore.logout();
          await router.replace({ name: 'LoginPage' });
          return Promise.reject("unable to refresh token");
        }
        return api_client(originalRequest);
      }
      catch (error) {
        return Promise.reject(error);
      }
      finally{
        createRefreshTokenInterceptor()
      }
    })
}

createRefreshTokenInterceptor()

export default api(api_client);
