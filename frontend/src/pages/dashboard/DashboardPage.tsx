
import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  ArrowUpRight,
  CalendarDays,
  CheckCircle2,
  Clock3,
  ClipboardList,
  FileText,
  MoreHorizontal,
  Truck,
  Wrench,
} from "lucide-react";

import StatCard from "../../components/dashboard/StatCard";
import {
  getWorkOrders,
  type WorkOrder,
} from "../../services/workOrderService";
import {
  getAllVehicles,
  type Vehicle,
} from "../../services/vehicleService";
import { useAuthStore, type UserRole } from "../../stores/authStore";

type WorkOrderStatus =
  | "In Progress"
  | "Scheduled"
  | "Pending"
  | "Completed"
  | "Cancelled";

type VehicleCategory =
  | "available"
  | "inUse"
  | "maintenance"
  | "inactive"
  | "other";

function getStatusLabel(status: string): WorkOrderStatus {
  switch (status.toUpperCase()) {
    case "IN_PROGRESS":
      return "In Progress";

    case "ASSIGNED":
    case "SCHEDULED":
    case "TECHNICIAN_ACCEPTED":
      return "Scheduled";

    case "COMPLETED":
    case "CUSTOMER_CONFIRMED":
    case "CLOSED":
      return "Completed";

    case "CANCELLED":
      return "Cancelled";

    case "PENDING":
    default:
      return "Pending";
  }
}

function getVehicleCategory(vehicle: Vehicle): VehicleCategory {
  if (!vehicle.is_active) return "inactive";

  switch (vehicle.status.toUpperCase()) {
    case "AVAILABLE":
      return "available";

    case "IN_USE":
      return "inUse";

    case "MAINTENANCE":
      return "maintenance";

    default:
      return "other";
  }
}

export default function DashboardPage() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);

  const canViewFleet =
    user?.role === "ADMIN" || user?.role === "DISPATCHER";

  const [workOrders, setWorkOrders] = useState<WorkOrder[]>([]);
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);

  const [isWorkOrdersLoading, setIsWorkOrdersLoading] = useState(true);
  const [isVehiclesLoading, setIsVehiclesLoading] = useState(true);

  const [workOrdersError, setWorkOrdersError] = useState("");
  const [vehiclesError, setVehiclesError] = useState("");

  const quickActions: {
    title: string;
    description: string;
    icon: typeof ClipboardList;
    path: string;
    roles: UserRole[];
  }[] = [
    {
      title: "View Work Orders",
      description: "Monitor and manage work orders",
      icon: ClipboardList,
      path: "/work-orders",
      roles: ["ADMIN", "DISPATCHER", "TECHNICIAN", "CUSTOMER"],
    },
    {
      title: "Review Requests",
      description: "Review incoming service requests",
      icon: FileText,
      path: "/service-requests",
      roles: ["ADMIN", "DISPATCHER", "CUSTOMER"],
    },
    {
      title: "Manage Technicians",
      description: "View technician information",
      icon: Wrench,
      path: "/technicians",
      roles: ["ADMIN", "DISPATCHER"],
    },
    {
      title: "Manage Vehicles",
      description: "View and manage fleet vehicles",
      icon: Truck,
      path: "/vehicles",
      roles: ["ADMIN", "DISPATCHER"],
    },
  ];

  const visibleQuickActions = quickActions.filter(
    (action) => user && action.roles.includes(user.role),
  );

  useEffect(() => {
    let isMounted = true;

    const loadDashboard = async () => {
      // Work Orders tersedia untuk semua role.
      try {
        const result = await getWorkOrders();

        if (!isMounted) return;

        setWorkOrders(result);
        setWorkOrdersError("");
      } catch (error) {
        if (!isMounted) return;

        setWorkOrdersError(
          error instanceof Error
            ? error.message
            : "Failed to load work orders.",
        );
      } finally {
        if (isMounted) {
          setIsWorkOrdersLoading(false);
        }
      }

      // Fleet hanya tersedia untuk ADMIN dan DISPATCHER.
      if (!canViewFleet) {
        if (isMounted) {
          setIsVehiclesLoading(false);
          setVehicles([]);
          setVehiclesError("");
        }

        return;
      }

      try {
        const result = await getAllVehicles();

        if (!isMounted) return;

        setVehicles(result);
        setVehiclesError("");
      } catch (error) {
        if (!isMounted) return;

        setVehiclesError(
          error instanceof Error
            ? error.message
            : "Failed to load vehicles.",
        );
      } finally {
        if (isMounted) {
          setIsVehiclesLoading(false);
        }
      }
    };

    void loadDashboard();

    return () => {
      isMounted = false;
    };
  }, [canViewFleet]);

  const sortedWorkOrders = useMemo(() => {
    return [...workOrders].sort(
      (a, b) =>
        new Date(b.created_at).getTime() -
        new Date(a.created_at).getTime(),
    );
  }, [workOrders]);

  const recentWorkOrders = sortedWorkOrders.slice(0, 5);

  const stats = useMemo(() => {
    return workOrders.reduce(
      (result, order) => {
        const status = getStatusLabel(order.status);

        if (status === "In Progress") result.inProgress += 1;
        if (status === "Scheduled") result.scheduled += 1;
        if (status === "Pending") result.pending += 1;
        if (status === "Completed") result.completed += 1;
        if (status === "Cancelled") result.cancelled += 1;

        return result;
      },
      {
        inProgress: 0,
        scheduled: 0,
        pending: 0,
        completed: 0,
        cancelled: 0,
      },
    );
  }, [workOrders]);

  const chartTotal =
    stats.inProgress +
    stats.scheduled +
    stats.pending +
    stats.completed +
    stats.cancelled;

  const chartStyle = useMemo(() => {
    if (chartTotal === 0) {
      return {
        background: "conic-gradient(#e5e7eb 0% 100%)",
      };
    }

    const segments = [
      { value: stats.inProgress, color: "#3b82f6" },
      { value: stats.scheduled, color: "#8b5cf6" },
      { value: stats.pending, color: "#f59e0b" },
      { value: stats.completed, color: "#22c55e" },
      { value: stats.cancelled, color: "#dc4c64" },
    ];

    let currentPosition = 0;

    const gradientStops = segments.map((segment) => {
      const start = currentPosition;
      currentPosition += (segment.value / chartTotal) * 100;

      return `${segment.color} ${start}% ${currentPosition}%`;
    });

    return {
      background: `conic-gradient(${gradientStops.join(", ")})`,
    };
  }, [chartTotal, stats]);

  const vehicleStats = useMemo(() => {
    return vehicles.reduce(
      (result, vehicle) => {
        const category = getVehicleCategory(vehicle);

        result[category] += 1;

        return result;
      },
      {
        available: 0,
        inUse: 0,
        maintenance: 0,
        inactive: 0,
        other: 0,
      },
    );
  }, [vehicles]);

  return (
    <div className="dashboard-page">
      {/* Dashboard Heading */}
      <section className="dashboard-heading">
        <div>
          <p className="dashboard-eyebrow">OVERVIEW</p>

          <h1 className="page-title">Dashboard</h1>

          <p className="page-description">
            Here&apos;s what&apos;s happening with your fleet and service
            operations today.
          </p>
        </div>

        <div className="dashboard-date-button">
          <CalendarDays size={17} />

          <span>
            {new Date().toLocaleDateString("en-GB", {
              day: "2-digit",
              month: "long",
              year: "numeric",
            })}
          </span>
        </div>
      </section>

      {/* Work Order Error */}
      {workOrdersError && (
        <div role="alert" className="dashboard-error">
          <p>Failed to load dashboard work orders.</p>
          <p>{workOrdersError}</p>

          <button
            type="button"
            onClick={() => window.location.reload()}
          >
            Try again
          </button>
        </div>
      )}

      {/* Statistics */}
      <section
        className="dashboard-stats"
        aria-label="Operational statistics"
      >
        <StatCard
          title="Total Work Orders"
          value={
            isWorkOrdersLoading ? "—" : String(workOrders.length)
          }
          description="All retrieved work orders"
          icon={ClipboardList}
          tone="blue"
        />

        <StatCard
          title="In Progress"
          value={
            isWorkOrdersLoading ? "—" : String(stats.inProgress)
          }
          description={
            <>
              <ArrowUpRight size={14} />
              Active work orders
            </>
          }
          icon={Wrench}
          tone="orange"
        />

        <StatCard
          title="Pending Work Orders"
          value={
            isWorkOrdersLoading ? "—" : String(stats.pending)
          }
          description={
            <>
              <Clock3 size={14} />
              Awaiting processing
            </>
          }
          icon={Clock3}
          tone="purple"
        />

        <StatCard
          title="Completed"
          value={
            isWorkOrdersLoading ? "—" : String(stats.completed)
          }
          description={
            <>
              <CheckCircle2 size={14} />
              Completed work orders
            </>
          }
          icon={CheckCircle2}
          tone="green"
        />
      </section>

      {/* Middle Dashboard */}
      <section className="dashboard-middle-grid">
        {/* Work Order Overview */}
        <article className="dashboard-panel workload-panel">
          <div className="panel-heading">
            <div>
              <h2 className="panel-title">Work Order Overview</h2>

              <p className="panel-description">
                Current distribution of work orders
              </p>
            </div>

            <button
              className="panel-more-button"
              type="button"
              aria-label="More work order overview options"
            >
              <MoreHorizontal size={20} />
            </button>
          </div>

          <div className="workload-content">
            <div
              className="workload-chart"
              style={chartStyle}
              role="img"
              aria-label={`Work order distribution: ${stats.inProgress} in progress, ${stats.scheduled} scheduled, ${stats.pending} pending, ${stats.completed} completed, ${stats.cancelled} cancelled`}
            >
              <div className="workload-chart-center">
                <span>Total orders</span>

                <strong>
                  {isWorkOrdersLoading ? "—" : chartTotal}
                </strong>

                <small>Retrieved orders</small>
              </div>
            </div>

            <div className="workload-legend">
              <div className="legend-item">
                <span className="legend-dot legend-dot--blue" />
                <span>In Progress</span>
                <strong>
                  {isWorkOrdersLoading ? "—" : stats.inProgress}
                </strong>
              </div>

              <div className="legend-item">
                <span className="legend-dot legend-dot--purple" />
                <span>Scheduled</span>
                <strong>
                  {isWorkOrdersLoading ? "—" : stats.scheduled}
                </strong>
              </div>

              <div className="legend-item">
                <span className="legend-dot legend-dot--orange" />
                <span>Pending</span>
                <strong>
                  {isWorkOrdersLoading ? "—" : stats.pending}
                </strong>
              </div>

              <div className="legend-item">
                <span className="legend-dot legend-dot--green" />
                <span>Completed</span>
                <strong>
                  {isWorkOrdersLoading ? "—" : stats.completed}
                </strong>
              </div>

              <div className="legend-item">
                <span className="legend-dot legend-dot--red" />
                <span>Cancelled</span>
                <strong>
                  {isWorkOrdersLoading ? "—" : stats.cancelled}
                </strong>
              </div>
            </div>
          </div>

          {stats.cancelled > 0 && !isWorkOrdersLoading && (
            <p className="panel-description">
              Cancelled work orders: {stats.cancelled}
            </p>
          )}
        </article>

        {/* Fleet Availability */}
        <article className="dashboard-panel fleet-panel">
          <div className="panel-heading">
            <div>
              <h2 className="panel-title">Fleet Availability</h2>

              <p className="panel-description">
                Vehicle status overview
              </p>
            </div>

            <Truck size={20} className="panel-heading-icon" />
          </div>

          {!canViewFleet ? (
            <div className="fleet-total">
              <span>
                Fleet overview is available to administrators and
                dispatchers.
              </span>
            </div>
          ) : (
            <>
              <div className="fleet-total">
                <strong>
                  {isVehiclesLoading ? "—" : vehicles.length}
                </strong>

                <span>
                  {vehiclesError
                    ? "Unable to retrieve vehicles"
                    : "Total vehicles"}
                </span>
              </div>

              {vehiclesError ? (
                <p className="panel-description" role="alert">
                  {vehiclesError}
                </p>
              ) : (
                <div className="fleet-status-list">
                  <div className="legend-item">
                    <span className="legend-dot legend-dot--green" />
                    <span>Available</span>
                    <strong>
                      {isVehiclesLoading ? "—" : vehicleStats.available}
                    </strong>
                  </div>

                  <div className="legend-item">
                    <span className="legend-dot legend-dot--blue" />
                    <span>In Use</span>
                    <strong>
                      {isVehiclesLoading ? "—" : vehicleStats.inUse}
                    </strong>
                  </div>

                  <div className="legend-item">
                    <span className="legend-dot legend-dot--orange" />
                    <span>Maintenance</span>
                    <strong>
                      {isVehiclesLoading ? "—" : vehicleStats.maintenance}
                    </strong>
                  </div>

                  <div className="legend-item">
                    <span className="legend-dot legend-dot--red" />
                    <span>Inactive</span>
                    <strong>
                      {isVehiclesLoading ? "—" : vehicleStats.inactive}
                    </strong>
                  </div>

                  {vehicleStats.other > 0 && (
                    <div className="legend-item">
                      <span className="legend-dot" />
                      <span>Other</span>
                      <strong>{vehicleStats.other}</strong>
                    </div>
                  )}
                </div>
              )}
            </>
          )}
        </article>
      </section>

      {/* Recent Work Orders */}
      <section className="dashboard-section recent-work-orders">
        <div className="section-header">
          <div>
            <h2>Recent Work Orders</h2>

            <p>Latest service jobs across your operations</p>
          </div>

          <button
            type="button"
            className="view-all-button"
            onClick={() => navigate("/work-orders")}
          >
            View All
          </button>
        </div>

        {isWorkOrdersLoading ? (
          <div className="dashboard-empty-state">
            Loading work orders...
          </div>
        ) : workOrdersError ? (
          <div className="dashboard-empty-state">
            {workOrdersError}
          </div>
        ) : workOrders.length === 0 ? (
          <div className="dashboard-empty-state">
            No work orders found.
          </div>
        ) : (
          <div className="recent-work-orders-table-wrapper">
            <table className="recent-work-orders-table">
              <thead>
                <tr>
                  <th>Work Order</th>
                  <th>Customer</th>
                  <th>Service Address</th>
                  <th>Status</th>
                  <th>Created</th>
                </tr>
              </thead>

              <tbody>
                {recentWorkOrders.map((workOrder) => (
                  <tr key={workOrder.id}>
                    <td>
                      <div className="work-order-info">
                        <strong>
                          WO-{String(workOrder.id).padStart(4, "0")}
                        </strong>

                        <span>{workOrder.title}</span>
                      </div>
                    </td>

                    <td>
                      Customer #{workOrder.customer_id}
                    </td>

                    <td>
                      <span
                        className="work-order-address"
                        title={workOrder.service_address}
                      >
                        {workOrder.service_address || "-"}
                      </span>
                    </td>

                    <td>
                      <span
                        className={`work-order-status status-${workOrder.status
                          .toLowerCase()
                          .replaceAll("_", "-")}`}
                      >
                        {workOrder.status
                          .replaceAll("_", " ")
                          .toLowerCase()
                          .replace(/\b\w/g, (char) =>
                            char.toUpperCase()
                          )}
                      </span>
                    </td>

                    <td>
                      {new Date(
                        workOrder.created_at
                      ).toLocaleDateString("en-GB", {
                        day: "2-digit",
                        month: "short",
                        year: "numeric",
                      })}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {/* Quick Actions */}
      <section className="dashboard-quick-actions">
        <div>
          <h2 className="panel-title">Quick Actions</h2>

          <p className="panel-description">
            Common tasks for your operations team
          </p>
        </div>

        <div className="quick-actions-grid">
          {visibleQuickActions.map((action) => {
            const Icon = action.icon;

            return (
              <button
                key={action.title}
                type="button"
                className="quick-action-card"
                onClick={() => navigate(action.path)}
              >
                <div className="quick-action-icon">
                  <Icon size={20} />
                </div>

                <div className="quick-action-content">
                  <h3>{action.title}</h3>
                  <p>{action.description}</p>
                </div>

                <ArrowUpRight size={18} />
              </button>
            );
          })}
        </div>
      </section>
    </div>
  );
}
