import axios from 'axios';

interface IApiResponse<T> {
  data: T;
  message: string;
}

const api_client = axios.create({
  baseURL: import.meta.env["VITE_API_BASE_URL"],
  timeout: 10000,
  withCredentials: true,
})


export {api_client, type IApiResponse};
