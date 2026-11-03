
import { useEffect, useState } from "react";
import {
  ArrowRight,
  Bell,
  CheckCircle2,
  ClipboardList,
  Clock3,
  Wrench,
} from "lucide-react";

import { workOrderApi, type WorkOrder } from "../api/workOrders";
import {
  notificationApi,
  type Notification,
} from "../api/notifications";

const STATUS_COLORS = [
  { dot: "dot-blue", fill: "fill-blue" },
  { dot: "dot-violet", fill: "fill-violet" },
  { dot: "dot-orange", fill: "fill-orange" },
  { dot: "dot-gray", fill: "fill-gray" },
];

function formatStatus(status: string) {
  return status
    .toLowerCase()
    .replaceAll("_", " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function formatDate(value: string) {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "-";
  }

  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

function formatRelativeTime(value: string) {
  const timestamp = new Date(value).getTime();

  if (Number.isNaN(timestamp)) {
    return "-";
  }

  const diff = Math.max(0, Date.now() - timestamp);
  const minutes = Math.floor(diff / 60000);

  if (minutes < 1) return "Just now";
  if (minutes < 60) return `${minutes}m ago`;

  const hours = Math.floor(minutes / 60);

  if (hours < 24) return `${hours}h ago`;

  const days = Math.floor(hours / 24);

  if (days < 30) return `${days}d ago`;

  return formatDate(value);
}

function getStatusClass(status: string) {
  return `status-${status.toLowerCase().replaceAll("_", "-")}`;
}

function getWorkOrderCounts(orders: WorkOrder[]) {
  return orders.reduce<Record<string, number>>((counts, order) => {
    const status = order.status.toUpperCase();
    counts[status] = (counts[status] ?? 0) + 1;
    return counts;
  }, {});
}

export default function DashboardPage() {
  const [workOrders, setWorkOrders] = useState<WorkOrder[]>([]);
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const [loadingOrders, setLoadingOrders] = useState(true);
  const [loadingNotifications, setLoadingNotifications] = useState(true);

  const [orderError, setOrderError] = useState("");
  const [notificationError, setNotificationError] = useState("");

  useEffect(() => {
    let active = true;

    async function loadDashboard() {
      const [ordersResult, notificationsResult] = await Promise.allSettled([
        workOrderApi.getAll(),
        notificationApi.getAll(),
      ]);

      if (!active) return;

      if (ordersResult.status === "fulfilled") {
        setWorkOrders(ordersResult.value);
        setOrderError("");
      } else {
        setOrderError(
          ordersResult.reason instanceof Error
            ? ordersResult.reason.message
            : "Failed to load work orders.",
        );
      }

      if (notificationsResult.status === "fulfilled") {
        setNotifications(notificationsResult.value);
        setNotificationError("");
      } else {
        setNotificationError(
          notificationsResult.reason instanceof Error
            ? notificationsResult.reason.message
            : "Failed to load notifications.",
        );
      }

      setLoadingOrders(false);
      setLoadingNotifications(false);
    }

    void loadDashboard();

    return () => {
      active = false;
    };
  }, []);

  const statusCounts = getWorkOrderCounts(workOrders);

  const pendingCount = statusCounts.PENDING ?? 0;
  const assignedCount = statusCounts.ASSIGNED ?? 0;
  const unreadCount = notifications.filter(
    (notification) => !notification.is_read,
  ).length;

  const workload = Object.entries(statusCounts).sort(
    ([statusA], [statusB]) => {
      const order = [
        "PENDING",
        "ASSIGNED",
        "SCHEDULED",
        "TECHNICIAN_ACCEPTED",
        "IN_PROGRESS",
        "COMPLETED",
        "CUSTOMER_CONFIRMED",
        "CLOSED",
        "CANCELLED",
      ];

      const indexA = order.indexOf(statusA);
      const indexB = order.indexOf(statusB);

      return (
        (indexA === -1 ? order.length : indexA) -
        (indexB === -1 ? order.length : indexB)
      );
    },
  );

  const stats = [
    {
      label: "Total Work Orders",
      value: workOrders.length,
      description: "Orders returned by API",
      icon: ClipboardList,
      color: "blue",
    },
    {
      label: "Pending Work Orders",
      value: pendingCount,
      description: "Awaiting processing",
      icon: Clock3,
      color: "violet",
    },
    {
      label: "Assigned Work Orders",
      value: assignedCount,
      description: "Assigned to technicians",
      icon: Wrench,
      color: "green",
    },
    {
      label: "Unread Notifications",
      value: unreadCount,
      description: "Unread notifications returned by API",
      icon: Bell,
      color: "orange",
    },
  ];

  const recentOrders = [...workOrders]
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() -
        new Date(a.created_at).getTime(),
    )
    .slice(0, 5);

  const recentNotifications = [...notifications]
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() -
        new Date(a.created_at).getTime(),
    )
    .slice(0, 4);

  return (
    <div className="dashboard">
      <div className="page-heading">
        <div>
          <div className="page-eyebrow">OVERVIEW</div>
          <h1>Dashboard</h1>
          <p>Here’s what’s happening across your service operations.</p>
        </div>

        <div className="heading-date">
          <Clock3 size={16} />
          <span>
            {new Date().toLocaleDateString("en-US", {
              month: "short",
              day: "numeric",
              year: "numeric",
            })}
          </span>
        </div>
      </div>

      {orderError && (
        <div className="dashboard-error" role="alert">
          <p>Unable to load work orders: {orderError}</p>
          <button onClick={() => window.location.reload()}>
            Retry
          </button>
        </div>
      )}

      <section className="stats-grid">
        {stats.map((stat) => {
          const Icon = stat.icon;

          const isLoading =
            stat.label === "Unread Notifications"
              ? loadingNotifications
              : loadingOrders;

          return (
            <article className="stat-card" key={stat.label}>
              <div className="stat-card-top">
                <span>{stat.label}</span>
                <div className={`stat-icon stat-icon-${stat.color}`}>
                  <Icon size={20} />
                </div>
              </div>

              <div className="stat-value">
                {isLoading ? "—" : stat.value.toLocaleString()}
              </div>

              <div className="stat-foot">
                <span>{stat.description}</span>
              </div>
            </article>
          );
        })}
      </section>

      <section className="dashboard-main-grid">
        <article className="panel workload-panel">
          <div className="panel-heading">
            <div>
              <h2>Work Order Overview</h2>
              <p>Current workload distribution</p>
            </div>
            <span className="period-select">All returned records</span>
          </div>

          <div className="workload-total">
            <strong>
              {loadingOrders ? "—" : workOrders.length.toLocaleString()}
            </strong>
            <span>work orders returned by API</span>
          </div>

          {loadingOrders ? (
            <p>Loading work order statistics...</p>
          ) : orderError ? (
            <p>Work order statistics are unavailable.</p>
          ) : workOrders.length === 0 ? (
            <p>No work orders found.</p>
          ) : (
            <div className="workload-bars">
              {workload.map(([status, count], index) => {
                const color =
                  STATUS_COLORS[index % STATUS_COLORS.length];

                const percentage =
                  workOrders.length > 0
                    ? (count / workOrders.length) * 100
                    : 0;

                return (
                  <div className="workload-row" key={status}>
                    <div className="workload-label">
                      <span>
                        <i className={`legend-dot ${color.dot}`} />
                        {formatStatus(status)}
                      </span>
                      <strong>{count}</strong>
                    </div>

                    <div className="progress-track">
                      <div
                        className={`progress-fill ${color.fill}`}
                        style={{ width: `${percentage}%` }}
                      />
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </article>

        <article className="panel activity-panel">
          <div className="panel-heading">
            <div>
              <h2>Service Activity</h2>
              <p>Latest notifications from the system</p>
            </div>

            <Bell size={18} />
          </div>

          {loadingNotifications ? (
            <p>Loading notifications...</p>
          ) : notificationError ? (
            <p role="alert">
              Unable to load notifications: {notificationError}
            </p>
          ) : recentNotifications.length === 0 ? (
            <p>No notifications yet.</p>
          ) : (
            <div className="activity-list">
              {recentNotifications.map((notification) => {
                const isCompleted =
                  notification.type === "WORK_COMPLETED";

                return (
                  <div className="activity-item" key={notification.id}>
                    <div
                      className={`activity-marker ${
                        isCompleted
                          ? "activity-marker-green"
                          : "activity-marker-blue"
                      }`}
                    >
                      {isCompleted ? (
                        <CheckCircle2 size={16} />
                      ) : (
                        <ClipboardList size={16} />
                      )}
                    </div>

                    <div className="activity-content">
                      <strong>{notification.title}</strong>
                      <p>{notification.message}</p>
                      <span>
                        {formatRelativeTime(notification.created_at)}
                        {!notification.is_read && " · Unread"}
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </article>
      </section>

      <section className="panel table-panel">
        <div className="panel-heading table-panel-heading">
          <div>
            <h2>Recent Work Orders</h2>
            <p>Latest service operations returned by the API</p>
          </div>

          <button
            className="text-button"
            onClick={() => window.location.assign("/work-orders")}
          >
            View all <ArrowRight size={16} />
          </button>
        </div>

        {loadingOrders ? (
          <p className="dashboard-table-message">
            Loading work orders...
          </p>
        ) : orderError ? (
          <p className="dashboard-table-message" role="alert">
            Work orders could not be loaded.
          </p>
        ) : recentOrders.length === 0 ? (
          <p className="dashboard-table-message">
            No work orders available.
          </p>
        ) : (
          <div className="table-scroll">
            <table className="work-orders-table">
              <thead>
                <tr>
                  <th>Work Order</th>
                  <th>Customer ID</th>
                  <th>Service Address</th>
                  <th>Created</th>
                  <th>Status</th>
                  <th />
                </tr>
              </thead>

              <tbody>
                {recentOrders.map((order) => (
                  <tr key={order.id}>
                    <td>
                      <div className="order-id">WO-{order.id}</div>
                      <div className="order-title">{order.title}</div>
                    </td>

                    <td>#{order.customer_id}</td>

                    <td>{order.service_address || "—"}</td>

                    <td>{formatDate(order.created_at)}</td>

                    <td>
                      <span
                        className={`status-badge ${getStatusClass(
                          order.status,
                        )}`}
                      >
                        <i />
                        {formatStatus(order.status)}
                      </span>
                    </td>

                    <td>
                      <button
                        className="row-action"
                        aria-label={`Open work order ${order.id}`}
                        onClick={() =>
                          window.location.assign(
                            `/work-orders/${order.id}`,
                          )
                        }
                      >
                        <ArrowRight size={17} />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <footer className="dashboard-footer">
        <span>SERVIXA Service Operations</span>
        <span>Fleet & Field Service Management</span>
      </footer>
    </div>
  );
}
