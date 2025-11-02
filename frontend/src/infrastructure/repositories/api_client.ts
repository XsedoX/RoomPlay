import axios, { type AxiosInstance } from 'axios';



const axiosParams = {
  baseURL: import.meta.env["VITE_API_BASE_URL"],
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    'Origin': import.meta.env["VITE_APP_ORIGIN"],
  },
  timeout: 10000,
  withCredentials: true,
};


const api_client = axios.create(axiosParams)

const api = (axiosInstance: AxiosInstance) => {
  return {
    get: <T>(url: string, config = {}) => axiosInstance.get<T>(url, config),
    post: <T>(url: string, body: unknown, config = {}) => axiosInstance.post<T>(url, body, config),
    put: <T>(url: string, body: unknown, config = {}) => axiosInstance.put<T>(url, body, config),
    delete: <T>(url: string, config = {}) => axiosInstance.delete<T>(url, config),
  };
};


export default api(api_client);
