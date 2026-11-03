
import { useEffect, useMemo, useState } from "react";
import {
  Bell,
  Check,
  CheckCheck,
  ClipboardList,
  Clock3,
  Filter,
  Package,
  RefreshCw,
  Search,
  Truck,
  UserRound,
  Wrench,
} from "lucide-react";
import {
  useNotificationStore,
  type NotificationItem,
} from "../../stores/notificationStore";

type NotificationCategory =
  | "WORK_ORDER"
  | "SERVICE_REQUEST"
  | "VEHICLE"
  | "INVENTORY"
  | "SYSTEM";

type ReadFilter = "ALL" | "UNREAD" | "READ";

const categoryLabels: Record<NotificationCategory, string> = {
  WORK_ORDER: "Work Orders",
  SERVICE_REQUEST: "Service Requests",
  VEHICLE: "Vehicles",
  INVENTORY: "Inventory",
  SYSTEM: "System",
};

function getCategory(
  notification: NotificationItem
): NotificationCategory {
  const value =
    `${notification.type} ${notification.reference}`.toUpperCase();

  if (value.includes("WORK_ORDER")) return "WORK_ORDER";
  if (value.includes("SERVICE_REQUEST")) return "SERVICE_REQUEST";
  if (value.includes("VEHICLE")) return "VEHICLE";

  if (
    value.includes("INVENTORY") ||
    value.includes("PART")
  ) {
    return "INVENTORY";
  }

  return "SYSTEM";
}

function getCategoryIcon(category: NotificationCategory) {
  switch (category) {
    case "WORK_ORDER":
      return ClipboardList;
    case "SERVICE_REQUEST":
      return Wrench;
    case "VEHICLE":
      return Truck;
    case "INVENTORY":
      return Package;
    default:
      return Bell;
  }
}

function formatNotificationDate(date: string) {
  const value = new Date(date);

  if (Number.isNaN(value.getTime())) return "-";

  return new Intl.DateTimeFormat("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(value);
}

export default function NotificationsPage() {
  const [readFilter, setReadFilter] =
    useState<ReadFilter>("ALL");

  const [categoryFilter, setCategoryFilter] =
    useState("ALL");

  const [search, setSearch] = useState("");
  const [actionLoading, setActionLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const notifications = useNotificationStore(
    (state) => state.notifications
  );

  const unreadCount = useNotificationStore(
    (state) => state.unreadCount
  );

  const isLoading = useNotificationStore(
    (state) => state.isLoading
  );

  const error = useNotificationStore(
    (state) => state.error
  );

  const fetchNotifications = useNotificationStore(
    (state) => state.fetchNotifications
  );

  const markAsRead = useNotificationStore(
    (state) => state.markAsRead
  );

  const markAllAsRead = useNotificationStore(
    (state) => state.markAllAsRead
  );

  // Initial load, auto-refresh, dan refresh saat tab aktif kembali.
  useEffect(() => {
    const refreshNotifications = () => {
      if (document.visibilityState === "visible") {
        void fetchNotifications();
      }
    };

    refreshNotifications();

    const interval = window.setInterval(
      refreshNotifications,
      30000
    );

    document.addEventListener(
      "visibilitychange",
      refreshNotifications
    );

    return () => {
      window.clearInterval(interval);

      document.removeEventListener(
        "visibilitychange",
        refreshNotifications
      );
    };
  }, [fetchNotifications]);

  // Kategori yang tersedia berdasarkan data notifikasi.
  const categories = useMemo(() => {
    return Array.from(
      new Set(
        notifications.map((item) => getCategory(item))
      )
    );
  }, [notifications]);

  // Search dan filter.
  const filteredNotifications = useMemo(() => {
    const keyword = search.trim().toLowerCase();

    return notifications.filter((notification) => {
      const category = getCategory(notification);

      const matchesSearch =
        !keyword ||
        notification.title.toLowerCase().includes(keyword) ||
        notification.message.toLowerCase().includes(keyword);

      const matchesRead =
        readFilter === "ALL" ||
        (readFilter === "UNREAD" && !notification.is_read) ||
        (readFilter === "READ" && notification.is_read);

      const matchesCategory =
        categoryFilter === "ALL" ||
        category === categoryFilter;

      return (
        matchesSearch &&
        matchesRead &&
        matchesCategory
      );
    });
  }, [
    notifications,
    readFilter,
    categoryFilter,
    search,
  ]);

  // Manual refresh.
  async function handleRefresh() {
    setRefreshing(true);

    try {
      await fetchNotifications();
    } finally {
      setRefreshing(false);
    }
  }

  // Mark one notification as read.
  async function handleMarkAsRead(id: number) {
    try {
      await markAsRead(id);
    } catch {
      window.alert(
        "Failed to mark notification as read."
      );
    }
  }

  // Mark all notifications as read.
  async function handleMarkAllAsRead() {
    if (unreadCount === 0 || actionLoading) return;

    setActionLoading(true);

    try {
      await markAllAsRead();
    } catch {
      window.alert(
        "Failed to mark all notifications as read."
      );
    } finally {
      setActionLoading(false);
    }
  }

  const readCount = notifications.filter(
    (item) => item.is_read
  ).length;

  return (
    <div className="notifications-page">
      {/* Page Heading */}
      <div className="page-heading">
        <div>
          <h1>Notifications</h1>
          <p>
            Stay updated on work orders, requests, and fleet operations.
          </p>
        </div>

        <div className="notifications-heading-actions">
          <button
            className="notifications-secondary-button"
            onClick={() => void handleRefresh()}
            disabled={isLoading || refreshing}
          >
            <RefreshCw
              size={16}
              className={
                isLoading || refreshing
                  ? "notification-spin"
                  : ""
              }
            />
            Refresh
          </button>

          <button
            className="notifications-primary-button"
            onClick={() => void handleMarkAllAsRead()}
            disabled={
              unreadCount === 0 ||
              actionLoading
            }
          >
            <CheckCheck size={17} />
            {actionLoading ? "Processing..." : "Mark All Read"}
          </button>
        </div>
      </div>

      {/* Statistics */}
      <div className="notifications-stats">
        <div className="notifications-stat-card">
          <div className="notifications-stat-icon blue">
            <Bell size={21} />
          </div>

          <div>
            <span>Total Notifications</span>
            <strong>{notifications.length}</strong>
            <small>All updates</small>
          </div>
        </div>

        <div className="notifications-stat-card">
          <div className="notifications-stat-icon orange">
            <Clock3 size={21} />
          </div>

          <div>
            <span>Unread</span>
            <strong>{unreadCount}</strong>
            <small>Require your attention</small>
          </div>
        </div>

        <div className="notifications-stat-card">
          <div className="notifications-stat-icon green">
            <CheckCheck size={21} />
          </div>

          <div>
            <span>Read</span>
            <strong>{readCount}</strong>
            <small>Already reviewed</small>
          </div>
        </div>
      </div>

      {/* Activity Center */}
      <section className="notifications-panel">
        <div className="notifications-panel-heading">
          <div>
            <h2>Activity Center</h2>
            <p>
              Your latest operational updates in one place.
            </p>
          </div>

          {unreadCount > 0 && (
            <span className="notifications-unread-count">
              {unreadCount} unread
            </span>
          )}
        </div>

        {/* Search and Category Filter */}
        <div className="notifications-toolbar">
          <label className="work-orders-search">
            <Search size={17} />

            <input
              type="search"
              placeholder="Search notifications..."
              value={search}
              onChange={(event) =>
                setSearch(event.target.value)
              }
            />
          </label>

          <label className="notifications-filter">
            <Filter size={16} />

            <select
              aria-label="Filter notifications by category"
              value={categoryFilter}
              onChange={(event) =>
                setCategoryFilter(event.target.value)
              }
            >
              <option value="ALL">
                All Categories
              </option>

              {categories.map((category) => (
                <option
                  value={category}
                  key={category}
                >
                  {categoryLabels[category]}
                </option>
              ))}
            </select>
          </label>
        </div>

        {/* Read Status Tabs */}
        <div
          className="notifications-tabs"
          role="group"
          aria-label="Read status"
        >
          <button
            className={
              readFilter === "ALL" ? "active" : ""
            }
            onClick={() => setReadFilter("ALL")}
          >
            All
            <span>{notifications.length}</span>
          </button>

          <button
            className={
              readFilter === "UNREAD" ? "active" : ""
            }
            onClick={() => setReadFilter("UNREAD")}
          >
            Unread
            <span>{unreadCount}</span>
          </button>

          <button
            className={
              readFilter === "READ" ? "active" : ""
            }
            onClick={() => setReadFilter("READ")}
          >
            Read
            <span>{readCount}</span>
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="notifications-error">
            {error}

            <button
              onClick={() => void handleRefresh()}
              disabled={isLoading || refreshing}
            >
              Try again
            </button>
          </div>
        )}

        {/* Notification List */}
        <div className="notifications-list">
          {isLoading && notifications.length === 0 ? (
            <div className="notifications-empty">
              <RefreshCw
                size={26}
                className="notification-spin"
              />

              <h3>Loading notifications</h3>

              <p>
                Please wait while we retrieve your updates.
              </p>
            </div>
          ) : (
            filteredNotifications.map((notification) => {
              const category = getCategory(notification);
              const CategoryIcon =
                getCategoryIcon(category);

              return (
                <article
                  key={notification.id}
                  className={`notification-item ${
                    notification.is_read
                      ? "is-read"
                      : "is-unread"
                  }`}
                >
                  <div
                    className={`notification-item-icon ${category.toLowerCase()}`}
                  >
                    <CategoryIcon size={20} />
                  </div>

                  <div className="notification-item-content">
                    <div className="notification-item-heading">
                      <h3>{notification.title}</h3>

                      {!notification.is_read && (
                        <span className="notification-new-indicator">
                          New
                        </span>
                      )}
                    </div>

                    <p>{notification.message}</p>

                    <div className="notification-item-meta">
                      <span>
                        {categoryLabels[category]}
                      </span>

                      <span className="notification-meta-dot">
                        ·
                      </span>

                      <span>
                        {formatNotificationDate(
                          notification.created_at
                        )}
                      </span>
                    </div>
                  </div>

                  <div className="notification-item-action">
                    {notification.is_read ? (
                      <span className="notification-read-label">
                        <Check size={14} />
                        Read
                      </span>
                    ) : (
                      <button
                        className="notification-mark-read"
                        onClick={() =>
                          void handleMarkAsRead(
                            notification.id
                          )
                        }
                      >
                        <Check size={15} />
                        Mark Read
                      </button>
                    )}
                  </div>
                </article>
              );
            })
          )}
        </div>

        {/* Empty State */}
        {!isLoading &&
          filteredNotifications.length === 0 &&
          !error && (
            <div className="notifications-empty">
              <UserRound size={30} />

              <h3>No notifications found</h3>

              <p>
                {notifications.length === 0
                  ? "You don't have any notifications yet."
                  : "Try changing your search keyword or filters."}
              </p>
            </div>
          )}

        {/* Results */}
        <div className="notifications-results">
          Showing {filteredNotifications.length} of{" "}
          {notifications.length} notifications
        </div>
      </section>
    </div>
  );
}
