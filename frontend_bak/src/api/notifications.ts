
import { api } from "./client";

export interface Notification {
  id: number;
  type: string;
  title: string;
  message: string;
  reference: string;
  is_read: boolean;
  read_at: string | null;
  created_at: string;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export const notificationApi = {
  async getAll(): Promise<Notification[]> {
    const response = await api.get<ApiResponse<Notification[]>>(
      "/notifications",
    );

    if (!response.data.success) {
      throw new Error(response.data.message || "Failed to load notifications");
    }

    return response.data.data;
  },
};
