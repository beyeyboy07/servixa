import { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
  ArrowLeft,
  CalendarDays,
  Clock3,
  CarFront,
  FileText,
  MapPin,
  RefreshCw,
  UserRound,
  UserPlus,
} from "lucide-react";
import Swal from "sweetalert2";

import {
  assignTechnician,
  reassignTechnician,
  assignVehicle,
  getWorkOrderById,
  type WorkOrder,
} from "../../services/workOrderService";

import {
  getAvailableVehicles,
  type Vehicle,
} from "../../services/vehicleService";

import {
  getAvailableTechnicians,
  type Technician,
} from "../../services/technicianService";

function formatWorkOrderId(id: number): string {
  return `WO-${String(id).padStart(4, "0")}`;
}

function formatStatus(status: string): string {
  return status
    .replaceAll("_", " ")
    .toLowerCase()
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function formatDateTime(value: string | null): string {
  if (!value) return "Not scheduled";

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "Not available";
  }

  return date.toLocaleString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function getErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof Error) {
    return err.message;
  }

  return fallback;
}

export default function WorkOrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [workOrder, setWorkOrder] = useState<WorkOrder | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  const [technicians, setTechnicians] = useState<Technician[]>([]);
  const [selectedTechnicianId, setSelectedTechnicianId] = useState("");
  const [isLoadingTechnicians, setIsLoadingTechnicians] = useState(true);
  const [technicianError, setTechnicianError] = useState("");

  const [isAssigning, setIsAssigning] = useState(false);
  const [assignmentError, setAssignmentError] = useState("");

  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [selectedVehicleId, setSelectedVehicleId] = useState("");
  const [isLoadingVehicles, setIsLoadingVehicles] = useState(true);
  const [vehicleError, setVehicleError] = useState("");
  const [isAssigningVehicle, setIsAssigningVehicle] = useState(false);

  const loadWorkOrder = useCallback(async () => {
    const workOrderId = Number(id);

    if (!id || !Number.isInteger(workOrderId) || workOrderId <= 0) {
      setError("Invalid work order ID.");
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError("");

      const data = await getWorkOrderById(workOrderId);
      setWorkOrder(data);
    } catch (err) {
      console.error("[SERVIXA] Failed to load work order detail:", err);
      setError(getErrorMessage(err, "Failed to retrieve work order."));
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  const loadTechnicians = useCallback(async () => {
    try {
      setIsLoadingTechnicians(true);
      setTechnicianError("");

      const data = await getAvailableTechnicians();
      setTechnicians(data);
    } catch (err) {
      console.error("[SERVIXA] Failed to load technicians:", err);
      setTechnicianError(
        getErrorMessage(err, "Failed to retrieve available technicians.")
      );
    } finally {
      setIsLoadingTechnicians(false);
    }
  }, []);

  const loadVehicles = useCallback(async () => {
    try {
      setIsLoadingVehicles(true);
      setVehicleError("");

      const data = await getAvailableVehicles();
      setVehicles(data);
    } catch (err) {
      console.error("[SERVIXA] Failed to load vehicles:", err);

      setVehicleError(
        getErrorMessage(err, "Failed to retrieve available vehicles.")
      );
    } finally {
      setIsLoadingVehicles(false);
    }
  }, []);

  useEffect(() => {
    void loadWorkOrder();
    void loadTechnicians();
    void loadVehicles();
  }, [loadWorkOrder, loadTechnicians, loadVehicles]);

  const handleTechnicianSubmit = async () => {
    const workOrderId = Number(id);
    const technicianId = Number(selectedTechnicianId);

    if (!Number.isInteger(workOrderId) || workOrderId <= 0) {
      setAssignmentError("Invalid work order ID.");
      return;
    }

    if (!Number.isInteger(technicianId) || technicianId <= 0) {
      setAssignmentError("Please select a technician.");
      return;
    }

    const isReassignment = Boolean(workOrder?.assignment);
    const selectedTechnician = technicians.find(
      (technician) => technician.id === technicianId
    );

    try {
      setIsAssigning(true);
      setAssignmentError("");

      if (isReassignment) {
        await reassignTechnician(workOrderId, technicianId);
      } else {
        await assignTechnician(workOrderId, technicianId);
      }

      await Swal.fire({
        icon: "success",
        title: isReassignment
          ? "Technician Reassigned!"
          : "Technician Assigned!",
        text: `${
          selectedTechnician?.name || `Technician #${technicianId}`
        } has been successfully ${isReassignment ? "reassigned" : "assigned"}.`,
        confirmButtonText: "OK",
        confirmButtonColor: "#2563eb",
      });

      setSelectedTechnicianId("");

      await Promise.all([loadWorkOrder(), loadTechnicians()]);
    } catch (err) {
      console.error("[SERVIXA] Failed to save technician assignment:", err);

      await Swal.fire({
        icon: "error",
        title: isReassignment ? "Reassignment Failed" : "Assignment Failed",
        text: getErrorMessage(
          err,
          isReassignment
            ? "Failed to reassign technician."
            : "Failed to assign technician."
        ),
        confirmButtonText: "OK",
        confirmButtonColor: "#dc2626",
      });
    } finally {
      setIsAssigning(false);
    }
  };

  const handleVehicleSubmit = async () => {
    const workOrderId = Number(id);
    const vehicleId = Number(selectedVehicleId);

    if (!Number.isInteger(workOrderId) || workOrderId <= 0) {
      setVehicleError("Invalid work order ID.");
      return;
    }

    if (!Number.isInteger(vehicleId) || vehicleId <= 0) {
      setVehicleError("Please select a vehicle.");
      return;
    }

    const selectedVehicle = vehicles.find(
      (vehicle) => vehicle.id === vehicleId
    );

    try {
      setIsAssigningVehicle(true);
      setVehicleError("");

      await assignVehicle(workOrderId, vehicleId);

      await Swal.fire({
        icon: "success",
        title: "Vehicle Assigned!",
        text: `${
          selectedVehicle?.plate_number || `Vehicle #${vehicleId}`
        } has been successfully assigned.`,
        confirmButtonText: "OK",
        confirmButtonColor: "#2563eb",
      });

      setSelectedVehicleId("");

      await Promise.all([loadWorkOrder(), loadVehicles()]);
    } catch (err) {
      console.error("[SERVIXA] Failed to assign vehicle:", err);

      await Swal.fire({
        icon: "error",
        title: "Vehicle Assignment Failed",
        text: getErrorMessage(err, "Failed to assign vehicle."),
        confirmButtonText: "OK",
        confirmButtonColor: "#dc2626",
      });
    } finally {
      setIsAssigningVehicle(false);
    }
  };

  if (isLoading) {
    return (
      <div className="work-order-detail-page">
        <div className="dashboard-panel work-order-detail-message">
          Loading work order details...
        </div>
      </div>
    );
  }

  if (error || !workOrder) {
    return (
      <div className="work-order-detail-page">
        <button
          type="button"
          className="work-order-back-button"
          onClick={() => navigate("/work-orders")}
        >
          <ArrowLeft size={17} />
          Back to Work Orders
        </button>

        <div className="dashboard-panel work-order-detail-message">
          <h2>Unable to load work order</h2>
          <p>{error || "Work order not found."}</p>

          <button
            type="button"
            className="work-order-retry-button"
            onClick={() => void loadWorkOrder()}
          >
            <RefreshCw size={16} />
            Try again
          </button>
        </div>
      </div>
    );
  }

  const currentAssignment = workOrder.assignment ?? null;
  const isAssigned = currentAssignment !== null;

  return (
    <div className="work-order-detail-page">
      <button
        type="button"
        className="work-order-back-button"
        onClick={() => navigate("/work-orders")}
      >
        <ArrowLeft size={17} />
        Back to Work Orders
      </button>

      <section className="work-order-detail-heading">
        <div>
          <p className="dashboard-eyebrow">OPERATIONS / WORK ORDERS / DETAIL</p>

          <h1 className="page-title">{formatWorkOrderId(workOrder.id)}</h1>

          <p className="page-description">
            Work order information and service progress.
          </p>
        </div>

        <span
          className={`status-badge status-badge--${workOrder.status
            .toLowerCase()
            .replaceAll("_", "-")}`}
        >
          <span className="status-dot" />
          {formatStatus(workOrder.status)}
        </span>
      </section>

      <div className="work-order-detail-grid">
        <section className="dashboard-panel work-order-detail-main">
          <div className="work-order-detail-section-heading">
            <FileText size={19} />
            <h2 className="panel-title">Service Information</h2>
          </div>

          <div className="work-order-detail-field">
            <span>Work Order Title</span>
            <strong>{workOrder.title || "—"}</strong>
          </div>

          <div className="work-order-detail-field">
            <span>Description</span>
            <p>{workOrder.description || "No description provided."}</p>
          </div>

          <div className="work-order-detail-field">
            <span>Service Address</span>
            <div className="work-order-detail-value">
              <MapPin size={17} />
              <strong>{workOrder.service_address || "—"}</strong>
            </div>
          </div>
        </section>

        <section className="dashboard-panel work-order-detail-side">
          <div className="work-order-detail-section-heading">
            <UserRound size={19} />
            <h2 className="panel-title">Related Information</h2>
          </div>

          <div className="work-order-detail-field">
            <span>Customer ID</span>
            <strong>#{workOrder.customer_id}</strong>
          </div>

          <div className="work-order-detail-field">
            <span>Service Request ID</span>
            <strong>#{workOrder.service_request_id}</strong>
          </div>
        </section>

        <section className="dashboard-panel work-order-detail-schedule">
          <div className="work-order-detail-section-heading">
            <CalendarDays size={19} />
            <h2 className="panel-title">Schedule</h2>
          </div>

          <div className="work-order-detail-schedule-grid">
            <div className="work-order-detail-field">
              <span>Scheduled Start</span>
              <div className="work-order-detail-value">
                <Clock3 size={16} />
                <strong>{formatDateTime(workOrder.scheduled_start)}</strong>
              </div>
            </div>

            <div className="work-order-detail-field">
              <span>Scheduled End</span>
              <div className="work-order-detail-value">
                <Clock3 size={16} />
                <strong>{formatDateTime(workOrder.scheduled_end)}</strong>
              </div>
            </div>
          </div>
        </section>

        <section className="dashboard-panel work-order-detail-metadata">
          <h2 className="panel-title">Record Information</h2>

          <div className="work-order-detail-field">
            <span>Created At</span>
            <strong>{formatDateTime(workOrder.created_at)}</strong>
          </div>

          <div className="work-order-detail-field">
            <span>Last Updated</span>
            <strong>{formatDateTime(workOrder.updated_at)}</strong>
          </div>
        </section>

        <section className="dashboard-panel work-order-detail-assignment">
          <div className="work-order-detail-section-heading">
            <UserPlus size={19} />
            <h2 className="panel-title">Technician Assignment</h2>
          </div>

          {isAssigned && (
            <div className="work-order-technician-info">
              <strong>Currently Assigned Technician</strong>
              <span>Technician ID: #{currentAssignment.technician_id}</span>
              <span>
                Assigned At: {formatDateTime(currentAssignment.assigned_at)}
              </span>
            </div>
          )}

          <div className="work-order-detail-field">
            <span>
              {isAssigned
                ? "Select Replacement Technician"
                : "Select Available Technician"}
            </span>

            <select
              className="work-order-technician-select"
              value={selectedTechnicianId}
              onChange={(event) => {
                setSelectedTechnicianId(event.target.value);
                setAssignmentError("");
              }}
              disabled={isLoadingTechnicians || isAssigning}
            >
              <option value="">
                {isLoadingTechnicians
                  ? "Loading technicians..."
                  : "Choose a technician"}
              </option>

              {technicians
                .filter(
                  (technician) =>
                    technician.id !== currentAssignment?.technician_id
                )
                .map((technician) => (
                  <option key={technician.id} value={technician.id}>
                    {technician.name} — {technician.specialization}
                  </option>
                ))}
            </select>
          </div>

          {selectedTechnicianId && (
            <div className="work-order-technician-info">
              {(() => {
                const technician = technicians.find(
                  (item) => item.id === Number(selectedTechnicianId)
                );

                if (!technician) return null;

                return (
                  <>
                    <strong>{technician.name}</strong>
                    <span>{technician.email}</span>
                    <span>{technician.phone}</span>
                    <span>{technician.specialization}</span>
                  </>
                );
              })()}
            </div>
          )}

          {technicianError && (
            <p className="assignment-error">{technicianError}</p>
          )}

          {assignmentError && (
            <p className="assignment-error">{assignmentError}</p>
          )}

          <button
            type="button"
            className="work-order-assign-button"
            onClick={() => void handleTechnicianSubmit()}
            disabled={
              isAssigning || isLoadingTechnicians || !selectedTechnicianId
            }
          >
            <UserPlus size={16} />
            {isAssigning
              ? isAssigned
                ? "Reassigning..."
                : "Assigning..."
              : isAssigned
              ? "Reassign Technician"
              : "Assign Technician"}
          </button>

          {!isLoadingTechnicians &&
            technicians.filter(
              (technician) => technician.id !== currentAssignment?.technician_id
            ).length === 0 &&
            !technicianError && (
              <p className="work-order-assignment-hint">
                No other available technicians found.
              </p>
            )}
        </section>

        <section className="dashboard-panel work-order-detail-assignment">
          <div className="work-order-detail-section-heading">
            <CarFront size={19} />
            <h2 className="panel-title">Vehicle Assignment</h2>
          </div>

          {workOrder.vehicle_assignment && (
            <div className="work-order-technician-info">
              <strong>Currently Assigned Vehicle</strong>

              <span>
                Vehicle ID: #{workOrder.vehicle_assignment.vehicle_id}
              </span>

              <span>
                Assigned At:{" "}
                {formatDateTime(workOrder.vehicle_assignment.assigned_at)}
              </span>
            </div>
          )}

          {!workOrder.vehicle_assignment && (
            <>
              <div className="work-order-detail-field">
                <span>Select Available Vehicle</span>

                <select
                  className="work-order-technician-select"
                  value={selectedVehicleId}
                  onChange={(event) => {
                    setSelectedVehicleId(event.target.value);
                    setVehicleError("");
                  }}
                  disabled={isLoadingVehicles || isAssigningVehicle}
                >
                  <option value="">
                    {isLoadingVehicles
                      ? "Loading vehicles..."
                      : "Choose a vehicle"}
                  </option>

                  {vehicles.map((vehicle) => (
                    <option key={vehicle.id} value={vehicle.id}>
                      {vehicle.plate_number} — {vehicle.brand} {vehicle.model} (
                      {vehicle.year})
                    </option>
                  ))}
                </select>
              </div>

              {selectedVehicleId && (
                <div className="work-order-technician-info">
                  {(() => {
                    const vehicle = vehicles.find(
                      (item) => item.id === Number(selectedVehicleId)
                    );

                    if (!vehicle) return null;

                    return (
                      <>
                        <strong>{vehicle.plate_number}</strong>
                        <span>
                          {vehicle.brand} {vehicle.model}
                        </span>
                        <span>Year: {vehicle.year}</span>
                        <span>Status: {vehicle.status}</span>
                      </>
                    );
                  })()}
                </div>
              )}

              {vehicleError && (
                <p className="assignment-error">{vehicleError}</p>
              )}

              <button
                type="button"
                className="work-order-assign-button"
                onClick={() => void handleVehicleSubmit()}
                disabled={
                  isAssigningVehicle || isLoadingVehicles || !selectedVehicleId
                }
              >
                <CarFront size={16} />

                {isAssigningVehicle ? "Assigning Vehicle..." : "Assign Vehicle"}
              </button>

              {!isLoadingVehicles && vehicles.length === 0 && !vehicleError && (
                <p className="work-order-assignment-hint">
                  No available vehicles found.
                </p>
              )}
            </>
          )}
        </section>
      </div>
    </div>
  );
}
