
import { useCallback, useEffect, useMemo, useState } from "react";
import {
  Ban,
  CalendarDays,
  CheckCircle2,
  ClipboardList,
  Clock3,
  Plus,
  Search,
  SlidersHorizontal,
  Wrench,
} from "lucide-react";

import Swal from "sweetalert2";
import { useNavigate } from "react-router-dom";

import {
  cancelWorkOrder,
  getWorkOrders,
  type WorkOrder,
} from "../../services/workOrderService";

import { useAuthStore } from "../../stores/authStore";

type WorkOrderStatus =
  | "PENDING"
  | "ASSIGNED"
  | "SCHEDULED"
  | "TECHNICIAN_ACCEPTED"
  | "IN_PROGRESS"
  | "COMPLETED"
  | "CUSTOMER_CONFIRMED"
  | "CLOSED"
  | "CANCELLED";

const statusLabels: Record<WorkOrderStatus, string> = {
  PENDING: "Pending",
  ASSIGNED: "Assigned",
  SCHEDULED: "Scheduled",
  TECHNICIAN_ACCEPTED: "Technician Accepted",
  IN_PROGRESS: "In Progress",
  COMPLETED: "Completed",
  CUSTOMER_CONFIRMED: "Customer Confirmed",
  CLOSED: "Closed",
  CANCELLED: "Cancelled",
};

function getStatusLabel(status: string): string {
  return (
    statusLabels[status as WorkOrderStatus] ??
    status.replaceAll("_", " ").toLowerCase().replace(/\b\w/g, (char) =>
      char.toUpperCase(),
    )
  );
}

function StatusBadge({ status }: { status: string }) {
  const className = status.toLowerCase().replaceAll("_", "-");

  return (
    <span className={`status-badge status-badge--${className}`}>
      <span className="status-dot" />
      {getStatusLabel(status)}
    </span>
  );
}

function formatWorkOrderId(id: number): string {
  return `WO-${String(id).padStart(4, "0")}`;
}

function formatDate(date: string | null): string {
  if (!date) {
    return "Not scheduled";
  }

  const parsedDate = new Date(date);

  if (Number.isNaN(parsedDate.getTime())) {
    return "Not scheduled";
  }

  return parsedDate.toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

export default function WorkOrdersPage() {
  const navigate = useNavigate();
  const { user } = useAuthStore();

  const canManageWorkOrders =
    user?.role === "ADMIN" || user?.role === "DISPATCHER";

  const [workOrders, setWorkOrders] = useState<WorkOrder[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("ALL");

  const [showCreateNotice, setShowCreateNotice] = useState(false);
  const [cancellingId, setCancellingId] = useState<number | null>(null);

  const loadWorkOrders = useCallback(async () => {
    try {
      setIsLoading(true);
      setError("");

      const data = await getWorkOrders();

      setWorkOrders(data);
    } catch (err) {
      console.error("[SERVIXA] Failed to load work orders:", err);

      setError(
        err instanceof Error
          ? err.message
          : "Failed to retrieve work orders.",
      );
    } finally {
      setIsLoading(false);
    }
  }, []);

  const handleCancelWorkOrder = async (order: WorkOrder) => {
    if (!canManageWorkOrders) {
      return;
    }

    const result = await Swal.fire({
      title: "Cancel Work Order?",
      html: `
        <p>Are you sure you want to cancel</p>
        <strong>${formatWorkOrderId(order.id)}</strong>
        <p style="margin-top: 12px; font-size: 13px; color: #64748b;">
          This action cannot be undone.
        </p>
      `,
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Yes, Cancel Work Order",
      cancelButtonText: "Keep Work Order",
      confirmButtonColor: "#dc2626",
      reverseButtons: true,
    });

    if (!result.isConfirmed) return;

    try {
      setCancellingId(order.id);

      const updatedWorkOrder = await cancelWorkOrder(order.id);

      setWorkOrders((current) =>
        current.map((item) =>
          item.id === order.id
            ? {
                ...item,
                ...updatedWorkOrder,
                status: "CANCELLED",
              }
            : item,
        ),
      );

      await Swal.fire({
        title: "Work Order Cancelled",
        text: `${formatWorkOrderId(order.id)} has been cancelled successfully.`,
        icon: "success",
        confirmButtonText: "OK",
      });
    } catch (err) {
      console.error("[SERVIXA] Failed to cancel work order:", err);

      const errorMessage =
        err instanceof Error
          ? err.message
          : "Failed to cancel work order.";

      await Swal.fire({
        title: "Cancellation Failed",
        text: errorMessage,
        icon: "error",
        confirmButtonText: "OK",
      });
    } finally {
      setCancellingId(null);
    }
  };

  useEffect(() => {
    void loadWorkOrders();
  }, [loadWorkOrders]);

  const filteredWorkOrders = useMemo(() => {
    const keyword = search.trim().toLowerCase();

    return workOrders.filter((order) => {
      const matchesSearch =
        !keyword ||
        formatWorkOrderId(order.id).toLowerCase().includes(keyword) ||
        String(order.id).includes(keyword) ||
        String(order.customer_id).includes(keyword) ||
        String(order.assignment?.technician_id ?? "").includes(keyword) ||
        String(order.vehicle_assignment?.vehicle_id ?? "").includes(keyword) ||
        order.title.toLowerCase().includes(keyword) ||
        (order.description ?? "").toLowerCase().includes(keyword) ||
        (order.service_address ?? "").toLowerCase().includes(keyword);

      const matchesStatus =
        statusFilter === "ALL" || order.status === statusFilter;

      return matchesSearch && matchesStatus;
    });
  }, [workOrders, search, statusFilter]);

  const counts = useMemo(() => {
    return {
      total: workOrders.length,

      active: workOrders.filter((order) =>
        [
          "ASSIGNED",
          "SCHEDULED",
          "TECHNICIAN_ACCEPTED",
          "IN_PROGRESS",
        ].includes(order.status),
      ).length,

      pending: workOrders.filter(
        (order) => order.status === "PENDING",
      ).length,

      completed: workOrders.filter((order) =>
        [
          "COMPLETED",
          "CUSTOMER_CONFIRMED",
          "CLOSED",
        ].includes(order.status),
      ).length,
    };
  }, [workOrders]);

  return (
    <div className="work-orders-page">
      <section className="work-orders-page-heading">
        <div>
          <p className="dashboard-eyebrow">OPERATIONS / WORK ORDERS</p>
          <h1 className="page-title">Work Orders</h1>
          <p className="page-description">
            Manage service jobs, technician assignments, and work progress.
          </p>
        </div>

        {canManageWorkOrders && (
          <button
            className="quick-action-button quick-action-button--primary"
            type="button"
            onClick={() => setShowCreateNotice((current) => !current)}
          >
            <Plus size={17} />
            New Work Order
          </button>
        )}
      </section>

      {canManageWorkOrders && showCreateNotice && (
        <div className="work-orders-notice">
          <div>
            <strong>Work order creation</strong>
            <p>
              Form pembuatan work order akan kita implementasikan
              setelah integrasi halaman ini selesai.
            </p>
          </div>

          <button
            type="button"
            className="work-orders-notice-close"
            onClick={() => setShowCreateNotice(false)}
            aria-label="Close notice"
          >
            ×
          </button>
        </div>
      )}

      <section className="work-orders-summary">
        <article className="work-orders-summary-card">
          <div className="work-orders-summary-icon summary-icon-blue">
            <ClipboardList size={20} />
          </div>
          <div>
            <p>Total Work Orders</p>
            <strong>{counts.total}</strong>
          </div>
        </article>

        <article className="work-orders-summary-card">
          <div className="work-orders-summary-icon summary-icon-orange">
            <Wrench size={20} />
          </div>
          <div>
            <p>Active Jobs</p>
            <strong>{counts.active}</strong>
          </div>
        </article>

        <article className="work-orders-summary-card">
          <div className="work-orders-summary-icon summary-icon-purple">
            <Clock3 size={20} />
          </div>
          <div>
            <p>Pending</p>
            <strong>{counts.pending}</strong>
          </div>
        </article>

        <article className="work-orders-summary-card">
          <div className="work-orders-summary-icon summary-icon-green">
            <CheckCircle2 size={20} />
          </div>
          <div>
            <p>Completed / Closed</p>
            <strong>{counts.completed}</strong>
          </div>
        </article>
      </section>

      <section className="dashboard-panel work-orders-list-panel">
        <div className="work-orders-list-heading">
          <div>
            <h2 className="panel-title">All Work Orders</h2>
            <p className="panel-description">
              Showing {filteredWorkOrders.length} of {workOrders.length} work
              orders
            </p>
          </div>
        </div>

        <div className="work-orders-toolbar">
          <label className="work-orders-search">
            <Search size={18} />
            <input
              type="search"
              placeholder="Search work order, customer, service..."
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </label>

          <label className="work-orders-filter">
            <SlidersHorizontal size={17} />
            <select
              value={statusFilter}
              onChange={(event) => setStatusFilter(event.target.value)}
              aria-label="Filter by work order status"
            >
              <option value="ALL">All statuses</option>

              {Object.entries(statusLabels).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </select>
          </label>
        </div>

        <div className="table-container">
          <table className="work-orders-table work-orders-full-table">
            <thead>
              <tr>
                <th>Work Order</th>
                <th>Customer / Service</th>
                <th>Technician</th>
                <th>Vehicle</th>
                <th>Scheduled Date</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>

            <tbody>
              {isLoading ? (
                <tr>
                  <td colSpan={7} className="work-orders-table-message">
                    Loading work orders...
                  </td>
                </tr>
              ) : error ? (
                <tr>
                  <td colSpan={7} className="work-orders-table-message">
                    <div className="work-orders-error">
                      <strong>Failed to load work orders</strong>
                      <p>{error}</p>
                      <button
                        type="button"
                        className="text-action"
                        onClick={() => void loadWorkOrders()}
                      >
                        Try again
                      </button>
                    </div>
                  </td>
                </tr>
              ) : (
                filteredWorkOrders.map((order) => (
                  <tr key={order.id}>
                    <td>
                      {canManageWorkOrders ? (
                        <button
                          type="button"
                          className="work-order-id work-order-detail-link"
                          onClick={() =>
                            navigate(`/work-orders/${order.id}`)
                          }
                        >
                          {formatWorkOrderId(order.id)}
                        </button>
                      ) : (
                        <span className="work-order-id">
                          {formatWorkOrderId(order.id)}
                        </span>
                      )}
                    </td>

                    <td>
                      <div className="table-primary-text">
                        Customer #{order.customer_id}
                      </div>
                      <div className="table-secondary-text">
                        {order.title}
                      </div>
                      {order.description && (
                        <div className="table-secondary-text">
                          {order.description}
                        </div>
                      )}
                    </td>

                    <td>
                      <span className="table-secondary-text">
                        {order.assignment
                          ? `Technician #${order.assignment.technician_id}`
                          : "Not assigned"}
                      </span>
                    </td>

                    <td>
                      <span className="table-secondary-text">
                        {order.vehicle_assignment
                          ? `Vehicle #${order.vehicle_assignment.vehicle_id}`
                          : "Not assigned"}
                      </span>
                    </td>

                    <td>
                      <span className="work-order-date">
                        <CalendarDays size={14} />
                        {formatDate(order.scheduled_start)}
                      </span>
                    </td>

                    <td>
                      <StatusBadge status={order.status} />
                    </td>

                    <td>
                      {canManageWorkOrders &&
                      ![
                        "COMPLETED",
                        "CUSTOMER_CONFIRMED",
                        "CLOSED",
                        "CANCELLED",
                      ].includes(order.status) ? (
                        <button
                          type="button"
                          className="work-order-cancel-button"
                          disabled={cancellingId === order.id}
                          onClick={() =>
                            void handleCancelWorkOrder(order)
                          }
                        >
                          <Ban size={15} />
                          {cancellingId === order.id
                            ? "Cancelling..."
                            : "Cancel"}
                        </button>
                      ) : (
                        <span className="table-secondary-text">—</span>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>

          {!isLoading && !error && filteredWorkOrders.length === 0 && (
            <div className="work-orders-empty">
              <Search size={26} />
              <strong>No work orders found</strong>
              <p>
                {workOrders.length === 0
                  ? "There are currently no work orders."
                  : "Try changing the search keyword or status filter."}
              </p>

              {(search || statusFilter !== "ALL") && (
                <button
                  type="button"
                  className="text-action"
                  onClick={() => {
                    setSearch("");
                    setStatusFilter("ALL");
                  }}
                >
                  Clear filters
                </button>
              )}
            </div>
          )}
        </div>

        <div className="work-orders-table-footer">
          <span>
            {filteredWorkOrders.length} result
            {filteredWorkOrders.length === 1 ? "" : "s"}
          </span>

          <span className="work-orders-demo-label">
            {isLoading
              ? "Loading..."
              : error
                ? "Connection error"
                : "Live data"}
          </span>
        </div>
      </section>
    </div>
  );
}
