import api from "./api";

export interface NotificationItem {
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

interface UnreadCount {
  count: number;
}

export const notificationService = {
  async getAll(): Promise<NotificationItem[]> {
    const response = await api.get<ApiResponse<NotificationItem[]>>(
      "/notifications",
    );

    return response.data.data ?? [];
  },

  async getUnread(): Promise<NotificationItem[]> {
    const response = await api.get<ApiResponse<NotificationItem[]>>(
      "/notifications/unread",
    );

    return response.data.data ?? [];
  },

  async getUnreadCount(): Promise<number> {
    const response = await api.get<ApiResponse<UnreadCount>>(
      "/notifications/unread/count",
    );

    return response.data.data?.count ?? 0;
  },

  async markAsRead(id: number): Promise<void> {
    await api.patch(`/notifications/${id}/read`);
  },

  async markAllAsRead(): Promise<void> {
    await api.patch("/notifications/read-all");
  },
};