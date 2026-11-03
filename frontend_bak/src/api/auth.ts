
import axios from "axios";
import { api } from "./client";
import type { AuthUser } from "../stores/authStore";

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface LoginResult {
  access_token: string;
  refresh_token: string;
  user: AuthUser;
}

export interface RefreshResult {
  access_token: string;
  refresh_token: string;
}

export const authApi = {
  async login(payload: LoginPayload): Promise<LoginResult> {
    const response = await api.post<ApiResponse<LoginResult>>(
      "/auth/login",
      payload,
    );

    return response.data.data;
  },

  async refresh(refreshToken: string): Promise<RefreshResult> {
    const response = await axios.post<ApiResponse<RefreshResult>>(
      `${import.meta.env.VITE_API_BASE_URL}/auth/refresh`,
      { refresh_token: refreshToken },
    );

    return response.data.data;
  },

  async logout(refreshToken: string): Promise<void> {
    await api.post("/auth/logout", {
      refresh_token: refreshToken,
    });
  },
};
