import axios, { AxiosError } from 'axios';
import { useAuthStore } from '@/entities/session/store';
import { toast } from '@/entities/notification/store';
import { API_URL } from '@/shared/config';

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (
    error: AxiosError<{ error?: string; message?: string; detail?: string }>
  ) => {
    // 1. Handle Auth Errors (401)
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
      toast.warn('Session Expired', 'Please sign in again.');
    }
    // 2. Handle Server Errors (5xx)
    else if (error.response && error.response.status >= 500) {
      toast.error(
        'Server Error',
        'Something went wrong on our end. Please try again later.'
      );
    }
    // 3. Handle Client Errors (4xx) - Validation, Business Logic
    else if (error.response && error.response.status >= 400) {
      const msg =
        error.response.data?.error ||
        error.response.data?.message ||
        error.response.data?.detail ||
        'Action failed';

      toast.error('Error', msg);
    }
    // 4. Network Errors (No response received)
    else if (error.code === 'ERR_NETWORK' || !error.response) {
      toast.error('Network Error', 'Please check your internet connection.');
    }

    return Promise.reject(error);
  }
);
