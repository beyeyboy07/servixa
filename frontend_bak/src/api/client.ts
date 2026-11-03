
import axios, {
  AxiosError,
  type InternalAxiosRequestConfig,
} from "axios";
import { useAuthStore } from "../stores/authStore";

const baseURL = import.meta.env.VITE_API_BASE_URL;

console.log("[SERVIXA] API Base URL:", baseURL);

export const api = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
});

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

interface RefreshResponse {
  access_token: string;
  refresh_token: string;
}

interface RetryConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

let refreshPromise: Promise<string> | null = null;

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken;

  if (token) {
    config.headers.set("Authorization", `Bearer ${token}`);
  }

  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as RetryConfig | undefined;

    if (
      error.response?.status !== 401 ||
      !originalRequest ||
      originalRequest._retry ||
      originalRequest.url?.includes("/auth/login") ||
      originalRequest.url?.includes("/auth/refresh") ||
      originalRequest.url?.includes("/auth/logout")
    ) {
      return Promise.reject(error);
    }

    const auth = useAuthStore.getState();

    if (!auth.refreshToken) {
      auth.clearAuth();
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    try {
      if (!refreshPromise) {
        const currentRefreshToken = auth.refreshToken;

        refreshPromise = axios
          .post<ApiResponse<RefreshResponse>>(
            `${baseURL}/auth/refresh`,
            { refresh_token: currentRefreshToken },
            { headers: { "Content-Type": "application/json" } },
          )
          .then((response) => {
            const data = response.data.data;

            useAuthStore.getState().setTokens(
              data.access_token,
              data.refresh_token,
            );

            return data.access_token;
          })
          .finally(() => {
            refreshPromise = null;
          });
      }

      const newAccessToken = await refreshPromise;

      originalRequest.headers.set(
        "Authorization",
        `Bearer ${newAccessToken}`,
      );

      return api(originalRequest);
    } catch (refreshError) {
      useAuthStore.getState().clearAuth();

      if (window.location.pathname !== "/login") {
        window.location.assign("/login");
      }

      return Promise.reject(refreshError);
    }
  },
);
