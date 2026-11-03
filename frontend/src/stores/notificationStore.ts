
import { create } from "zustand";
import {
  notificationService,
  type NotificationItem,
} from "../services/notificationService";

export type { NotificationItem } from "../services/notificationService";

interface NotificationState {
  notifications: NotificationItem[];
  unreadCount: number;
  isLoading: boolean;
  error: string | null;

  fetchNotifications: () => Promise<void>;
  fetchUnreadCount: () => Promise<void>;
  markAsRead: (id: number) => Promise<void>;
  markAllAsRead: () => Promise<void>;
  reset: () => void;
}

export const useNotificationStore = create<NotificationState>(
  (set, get) => ({
    notifications: [],
    unreadCount: 0,
    isLoading: false,
    error: null,

    fetchNotifications: async () => {
      set({ isLoading: true, error: null });

      try {
        const notifications = await notificationService.getAll();

        set({
          notifications,
          unreadCount: notifications.filter(
            (item) => !item.is_read
          ).length,
          isLoading: false,
          error: null,
        });
      } catch {
        set({
          error: "Failed to load notifications.",
          isLoading: false,
        });
      }
    },

    fetchUnreadCount: async () => {
      try {
        const unreadCount =
          await notificationService.getUnreadCount();

        set({ unreadCount });
      } catch {
        // Pertahankan unread count sebelumnya jika request gagal.
      }
    },

    markAsRead: async (id: number) => {
      const notification = get().notifications.find(
        (item) => item.id === id
      );

      if (notification?.is_read) return;

      try {
        await notificationService.markAsRead(id);

        set((state) => ({
          notifications: state.notifications.map((item) =>
            item.id === id
              ? {
                  ...item,
                  is_read: true,
                  read_at: new Date().toISOString(),
                }
              : item
          ),
          unreadCount:
            notification && !notification.is_read
              ? Math.max(0, state.unreadCount - 1)
              : state.unreadCount,
          error: null,
        }));

        // Sinkronkan jika notifikasi tidak ditemukan di state lokal.
        if (!notification) {
          await get().fetchUnreadCount();
        }
      } catch (error) {
        set({ error: "Failed to mark notification as read." });
        throw error;
      }
    },

    markAllAsRead: async () => {
      try {
        await notificationService.markAllAsRead();

        const readAt = new Date().toISOString();

        set((state) => ({
          notifications: state.notifications.map((item) => ({
            ...item,
            is_read: true,
            read_at: item.read_at ?? readAt,
          })),
          unreadCount: 0,
          error: null,
        }));
      } catch (error) {
        set({ error: "Failed to mark all notifications as read." });
        throw error;
      }
    },

    reset: () => {
      set({
        notifications: [],
        unreadCount: 0,
        isLoading: false,
        error: null,
      });
    },
  })
);
