
import { useCallback, useEffect, useMemo, useState, type FormEvent } from "react";
import {
  CarFront,
  CheckCircle2,
  Plus,
  Search,
  SlidersHorizontal,
  Truck,
  X,
  RefreshCw,
  Pencil,
  Power,
} from "lucide-react";
import Swal from "sweetalert2";

import {
  activateVehicle,
  createVehicle,
  deactivateVehicle,
  getAllVehicles,
  updateVehicle,
  type Vehicle,
} from "../../services/vehicleService";

import "./VehiclesPage.css";


type VehicleForm = {
  plate_number: string;
  brand: string;
  model: string;
  year: string;
  status: string;
};

const initialForm: VehicleForm = {
  plate_number: "",
  brand: "",
  model: "",
  year: "",
  status: "AVAILABLE",
};

function formatStatus(status: string): string {
  return status
    .replaceAll("_", " ")
    .toLowerCase()
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function formatYear(year: number): string {
  return String(year);
}

function getErrorMessage(err: unknown, fallback: string): string {
  return err instanceof Error ? err.message : fallback;
}

export default function VehiclesPage() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("ALL");

  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | null>(null);

  const [showAddModal, setShowAddModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);

  const [form, setForm] = useState<VehicleForm>(initialForm);
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [togglingVehicleId, setTogglingVehicleId] = useState<number | null>(
    null,
  );

  const loadVehicles = useCallback(async (): Promise<Vehicle[] | undefined> => {
    try {
      setIsLoading(true);
      setError("");

      const data = await getAllVehicles();
      setVehicles(data);

      return data;
    } catch (err) {
      console.error("[SERVIXA] Failed to load vehicles:", err);
      setError(getErrorMessage(err, "Failed to retrieve vehicles."));
      return undefined;
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadVehicles();
  }, [loadVehicles]);

  const filteredVehicles = useMemo(() => {
    const keyword = search.trim().toLowerCase();

    return vehicles.filter((vehicle) => {
      const matchesSearch =
        !keyword ||
        vehicle.plate_number.toLowerCase().includes(keyword) ||
        vehicle.brand.toLowerCase().includes(keyword) ||
        vehicle.model.toLowerCase().includes(keyword) ||
        String(vehicle.year).includes(keyword);

      let matchesStatus = true;

      if (statusFilter === "INACTIVE") {
        matchesStatus = !vehicle.is_active;
      } else if (statusFilter !== "ALL") {
        matchesStatus =
          vehicle.is_active && vehicle.status === statusFilter;
      }

      return matchesSearch && matchesStatus;
    });
  }, [vehicles, search, statusFilter]);

  const totalCount = vehicles.length;

  const availableCount = vehicles.filter(
    (vehicle) => vehicle.is_active && vehicle.status === "AVAILABLE",
  ).length;

  const inUseCount = vehicles.filter(
    (vehicle) => vehicle.is_active && vehicle.status === "IN_USE",
  ).length;

  const maintenanceCount = vehicles.filter(
    (vehicle) => vehicle.is_active && vehicle.status === "MAINTENANCE",
  ).length;

  const inactiveCount = vehicles.filter(
    (vehicle) => !vehicle.is_active,
  ).length;

  const handleFormChange = (field: keyof VehicleForm, value: string) => {
    setForm((current) => ({
      ...current,
      [field]: value,
    }));

    setFormError("");
  };

  const validateForm = () => {
    const plateNumber = form.plate_number.trim();
    const brand = form.brand.trim();
    const model = form.model.trim();
    const year = Number(form.year);

    if (!plateNumber || !brand || !model || !form.year) {
      setFormError("Please complete all vehicle information.");
      return false;
    }

    if (!Number.isInteger(year) || year <= 0) {
      setFormError("Please enter a valid vehicle year.");
      return false;
    }

    if (
      showEditModal &&
      !["AVAILABLE", "IN_USE", "MAINTENANCE"].includes(form.status)
    ) {
      setFormError("Please select a valid vehicle status.");
      return false;
    }

    return true;
  };

  const handleAddVehicle = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!validateForm()) return;

    const plateNumber = form.plate_number.trim();

    try {
      setIsSubmitting(true);
      setFormError("");

      await createVehicle({
        plate_number: plateNumber,
        brand: form.brand.trim(),
        model: form.model.trim(),
        year: Number(form.year),
      });

      setShowAddModal(false);
      setForm(initialForm);

      await Swal.fire({
        icon: "success",
        title: "Vehicle Added!",
        text: `${plateNumber.toUpperCase()} has been successfully registered.`,
        confirmButtonText: "OK",
        confirmButtonColor: "#2563eb",
      });

      await loadVehicles();
    } catch (err) {
      console.error("[SERVIXA] Failed to create vehicle:", err);
      setFormError(getErrorMessage(err, "Failed to register vehicle."));
    } finally {
      setIsSubmitting(false);
    }
  };

  const openEditModal = (vehicle: Vehicle) => {
    setSelectedVehicle(vehicle);

    setForm({
      plate_number: vehicle.plate_number,
      brand: vehicle.brand,
      model: vehicle.model,
      year: String(vehicle.year),
      status: vehicle.status,
    });

    setFormError("");
    setShowEditModal(true);
  };

  const handleUpdateVehicle = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!selectedVehicle || !validateForm()) return;

    try {
      setIsSubmitting(true);
      setFormError("");

      await updateVehicle(selectedVehicle.id, {
        plate_number: form.plate_number.trim(),
        brand: form.brand.trim(),
        model: form.model.trim(),
        year: Number(form.year),
        status: form.status,
      });

      setShowEditModal(false);
      setSelectedVehicle(null);
      setForm(initialForm);

      await Swal.fire({
        icon: "success",
        title: "Vehicle Updated!",
        text: "Vehicle information has been successfully updated.",
        confirmButtonText: "OK",
        confirmButtonColor: "#2563eb",
      });

      await loadVehicles();
    } catch (err) {
      console.error("[SERVIXA] Failed to update vehicle:", err);
      setFormError(getErrorMessage(err, "Failed to update vehicle."));
    } finally {
      setIsSubmitting(false);
    }
  };


  const handleToggleVehicleActive = async (vehicle: Vehicle) => {
    const shouldActivate = !vehicle.is_active;

    if (!shouldActivate && vehicle.status === "IN_USE") {
      await Swal.fire({
        icon: "warning",
        title: "Vehicle Currently In Use",
        text: "This vehicle cannot be deactivated while it is assigned to an active work order.",
        confirmButtonText: "OK",
        confirmButtonColor: "#2563eb",
      });

      return;
    }

    const result = await Swal.fire({
      icon: "warning",
      title: shouldActivate ? "Activate Vehicle?" : "Deactivate Vehicle?",
      text: shouldActivate
        ? `Are you sure you want to activate ${vehicle.plate_number}?`
        : `Are you sure you want to deactivate ${vehicle.plate_number}?`,
      showCancelButton: true,
      confirmButtonText: shouldActivate ? "Yes, Activate" : "Yes, Deactivate",
      cancelButtonText: "Cancel",
      confirmButtonColor: shouldActivate ? "#2563eb" : "#dc2626",
      reverseButtons: true,
    });

    if (!result.isConfirmed) return;

    try {
      setTogglingVehicleId(vehicle.id);

      if (shouldActivate) {
        await activateVehicle(vehicle.id);
      } else {
        await deactivateVehicle(vehicle.id);
      }

      const updatedVehicles = await loadVehicles();

      const updatedVehicle = updatedVehicles?.find(
        (item) => item.id === vehicle.id,
      );

      if (updatedVehicle) {
        setSelectedVehicle(updatedVehicle);
      }

      await Swal.fire({
        icon: "success",
        title: shouldActivate ? "Vehicle Activated!" : "Vehicle Deactivated!",
        text: shouldActivate
          ? `${vehicle.plate_number} is now active and available for fleet management.`
          : `${vehicle.plate_number} has been deactivated.`,
        confirmButtonText: "OK",
        confirmButtonColor: "#2563eb",
      });
    } catch (err) {
      console.error("[SERVIXA] Failed to change vehicle active status:", err);

      await Swal.fire({
        icon: "error",
        title: "Action Failed",
        text: getErrorMessage(
          err,
          `Failed to ${shouldActivate ? "activate" : "deactivate"} vehicle.`,
        ),
        confirmButtonText: "OK",
        confirmButtonColor: "#2563eb",
      });
    } finally {
      setTogglingVehicleId(null);
    }
  };

  const closeAddModal = () => {
    if (isSubmitting) return;

    setShowAddModal(false);
    setForm(initialForm);
    setFormError("");
  };

  const closeEditModal = () => {
    if (isSubmitting) return;

    setShowEditModal(false);
    setForm(initialForm);
    setFormError("");
  };

  const closeDetailModal = () => {
    setSelectedVehicle(null);
  };

  return (
    <div className="vehicles-page">
      <div className="page-heading">
        <div>
          <h1>Vehicles</h1>
          <p>Manage your fleet and monitor vehicle availability.</p>
        </div>

        <button
          type="button"
          className="vehicles-add-button"
          onClick={() => {
            setForm(initialForm);
            setFormError("");
            setShowAddModal(true);
          }}
        >
          <Plus size={18} />
          Add Vehicle
        </button>
      </div>

      <div className="vehicles-stats">
        <div className="vehicles-stat-card">
          <div className="vehicles-stat-icon blue">
            <Truck size={21} />
          </div>
          <div>
            <span>Total Fleet</span>
            <strong>{totalCount}</strong>
            <small>All registered vehicles</small>
          </div>
        </div>

        <div className="vehicles-stat-card">
          <div className="vehicles-stat-icon green">
            <CheckCircle2 size={21} />
          </div>
          <div>
            <span>Available</span>
            <strong>{availableCount}</strong>
            <small>Ready for assignment</small>
          </div>
        </div>

        <div className="vehicles-stat-card">
          <div className="vehicles-stat-icon orange">
            <Truck size={21} />
          </div>
          <div>
            <span>In Use</span>
            <strong>{inUseCount}</strong>
            <small>Currently assigned</small>
          </div>
        </div>

        <div className="vehicles-stat-card">
          <div className="vehicles-stat-icon purple">
            <CarFront size={21} />
          </div>
          <div>
            <span>Maintenance</span>
            <strong>{maintenanceCount}</strong>
            <small>Under maintenance</small>
          </div>
        </div>
      </div>

      <section className="vehicles-panel">
        <div className="vehicles-toolbar">
          <div>
            <h2>Fleet Overview</h2>
            <p>View and manage your registered vehicles.</p>
          </div>

          <div className="vehicles-toolbar-controls">
            <label className="work-orders-search">
              <Search size={17} />
              <input
                type="search"
                placeholder="Search vehicles..."
                value={search}
                onChange={(event) => setSearch(event.target.value)}
              />
            </label>

            <label className="work-orders-filter">
              <SlidersHorizontal size={16} />
              <select
                value={statusFilter}
                onChange={(event) => setStatusFilter(event.target.value)}
                aria-label="Filter vehicles by status"
              >
                <option value="ALL">All Status</option>
                <option value="AVAILABLE">Available</option>
                <option value="IN_USE">In Use</option>
                <option value="MAINTENANCE">Maintenance</option>
                <option value="INACTIVE">Inactive</option>
              </select>
            </label>

            <button
              type="button"
              className="work-order-retry-button"
              onClick={() => void loadVehicles()}
              disabled={isLoading}
              title="Refresh vehicles"
            >
              <RefreshCw size={16} />
            </button>
          </div>
        </div>

        {isLoading ? (
          <div className="vehicles-empty">Loading vehicles...</div>
        ) : error ? (
          <div className="vehicles-empty">
            <CarFront size={34} />
            <h3>Unable to load vehicles</h3>
            <p>{error}</p>

            <button
              type="button"
              className="work-order-retry-button"
              onClick={() => void loadVehicles()}
            >
              <RefreshCw size={16} />
              Try Again
            </button>
          </div>
        ) : (
          <>
            <div className="vehicles-grid">
              {filteredVehicles.map((vehicle) => (
                <article className="vehicle-card" key={vehicle.id}>
                  <div className="vehicle-card-top">
                    <div className="vehicle-icon">
                      <CarFront size={26} />
                    </div>

                    <span
                      className={`vehicle-status ${
                        vehicle.is_active
                          ? vehicle.status.toLowerCase()
                          : "inactive"
                      }`}
                    >
                      <span />
                      {vehicle.is_active
                        ? formatStatus(vehicle.status)
                        : "Inactive"}
                    </span>
                  </div>

                  <div className="vehicle-card-heading">
                    <h3>
                      {vehicle.brand} {vehicle.model}
                    </h3>
                    <p>{vehicle.year}</p>
                  </div>

                  <div className="vehicle-plate">
                    <span>LICENSE PLATE</span>
                    <strong>{vehicle.plate_number}</strong>
                  </div>

                  <div className="vehicle-card-footer">
                    <span>Vehicle #{vehicle.id}</span>

                    <button
                      type="button"
                      onClick={() => setSelectedVehicle(vehicle)}
                    >
                      View Details
                    </button>
                  </div>
                </article>
              ))}
            </div>

            {filteredVehicles.length === 0 && (
              <div className="vehicles-empty">
                <CarFront size={34} />
                <h3>No vehicles found</h3>
                <p>
                  Try another search keyword or change the status filter.
                </p>
              </div>
            )}
          </>
        )}

        <div className="vehicles-results">
          Showing {filteredVehicles.length} of {vehicles.length} vehicles
        </div>
      </section>

      {selectedVehicle && !showEditModal && (
        <div
          className="request-modal-backdrop"
          onClick={closeDetailModal}
        >
          <section
            className="request-modal vehicle-detail-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="vehicle-detail-title"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="request-modal-heading">
              <div>
                <h2 id="vehicle-detail-title">
                  {selectedVehicle.brand} {selectedVehicle.model}
                </h2>
                <p>Vehicle #{selectedVehicle.id}</p>
              </div>

              <button
                type="button"
                className="vehicle-close-button"
                onClick={closeDetailModal}
                aria-label="Close vehicle details"
              >
                <X size={20} />
              </button>
            </div>

            <div className="vehicle-detail-banner">
              <CarFront size={36} />
              <div>
                <strong>{selectedVehicle.plate_number}</strong>

                <span
                  className={`vehicle-status ${
                    selectedVehicle.is_active
                      ? selectedVehicle.status.toLowerCase()
                      : "inactive"
                  }`}
                >
                  <span />
                  {selectedVehicle.is_active
                    ? formatStatus(selectedVehicle.status)
                    : "Inactive"}
                </span>
              </div>
            </div>

            <div className="request-detail-grid">
              <div>
                <span>Brand</span>
                <strong>{selectedVehicle.brand}</strong>
              </div>

              <div>
                <span>Model</span>
                <strong>{selectedVehicle.model}</strong>
              </div>

              <div>
                <span>Year</span>
                <strong>{formatYear(selectedVehicle.year)}</strong>
              </div>

              <div>
                <span>Status</span>
                <strong>{formatStatus(selectedVehicle.status)}</strong>
              </div>

              <div>
                <span>Active</span>
                <strong>{selectedVehicle.is_active ? "Yes" : "No"}</strong>
              </div>

              <div>
                <span>Vehicle ID</span>
                <strong>#{selectedVehicle.id}</strong>
              </div>
            </div>

            <div className="vehicle-detail-actions">
              <button
                type="button"
                className="secondary-button"
                onClick={closeDetailModal}
                disabled={togglingVehicleId !== null}
              >
                Close
              </button>

              <button
                type="button"
                className={`vehicle-toggle-button ${
                  selectedVehicle.is_active ? "deactivate" : "activate"
                }`}
                onClick={() => void handleToggleVehicleActive(selectedVehicle)}
                disabled={
                  togglingVehicleId !== null ||
                  (selectedVehicle.is_active && selectedVehicle.status === "IN_USE")
                }
              >
                <Power size={16} />

                {togglingVehicleId === selectedVehicle.id
                  ? "Processing..."
                  : selectedVehicle.is_active
                    ? "Deactivate"
                    : "Activate"}
              </button>

              <button
                type="button"
                className="vehicles-add-button"
                onClick={() => openEditModal(selectedVehicle)}
                disabled={togglingVehicleId !== null}
              >
                <Pencil size={16} />
                Edit Vehicle
              </button>
            </div>
          </section>
        </div>
      )}

      {(showAddModal || showEditModal) && (
        <div
          className="request-modal-backdrop"
          onClick={showEditModal ? closeEditModal : closeAddModal}
        >
          <section
            className="request-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="vehicle-form-title"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="request-modal-heading">
              <div>
                <h2 id="vehicle-form-title">
                  {showEditModal ? "Edit Vehicle" : "Add Vehicle"}
                </h2>
                <p>
                  {showEditModal
                    ? "Update your vehicle information."
                    : "Register a new vehicle in your fleet."}
                </p>
              </div>

              <button
                type="button"
                className="vehicle-close-button"
                onClick={showEditModal ? closeEditModal : closeAddModal}
                aria-label="Close"
                disabled={isSubmitting}
              >
                <X size={20} />
              </button>
            </div>

            <form
              onSubmit={showEditModal ? handleUpdateVehicle : handleAddVehicle}
            >
              <div className="request-detail-grid vehicle-form-grid">
                <div className="vehicle-form-field">
                  <label htmlFor="vehicle-plate-number">Plate Number</label>
                  <input
                    id="vehicle-plate-number"
                    type="text"
                    placeholder="e.g. B 9123 TCA"
                    value={form.plate_number}
                    onChange={(event) =>
                      handleFormChange("plate_number", event.target.value)
                    }
                    disabled={isSubmitting}
                    required
                  />
                </div>

                <div className="vehicle-form-field">
                  <label htmlFor="vehicle-brand">Brand</label>
                  <input
                    id="vehicle-brand"
                    type="text"
                    placeholder="e.g. Toyota"
                    value={form.brand}
                    onChange={(event) =>
                      handleFormChange("brand", event.target.value)
                    }
                    disabled={isSubmitting}
                    required
                  />
                </div>

                <div className="vehicle-form-field">
                  <label htmlFor="vehicle-model">Model</label>
                  <input
                    id="vehicle-model"
                    type="text"
                    placeholder="e.g. Avanza"
                    value={form.model}
                    onChange={(event) =>
                      handleFormChange("model", event.target.value)
                    }
                    disabled={isSubmitting}
                    required
                  />
                </div>

                <div className="vehicle-form-field">
                  <label htmlFor="vehicle-year">Year</label>
                  <input
                    id="vehicle-year"
                    type="number"
                    placeholder="e.g. 2024"
                    min="1900"
                    max="2100"
                    value={form.year}
                    onChange={(event) =>
                      handleFormChange("year", event.target.value)
                    }
                    disabled={isSubmitting}
                    required
                  />
                </div>

                {showEditModal && (
                  <div className="vehicle-form-field">
                    <label htmlFor="vehicle-status">Vehicle Status</label>
                    <select
                      id="vehicle-status"
                      value={form.status}
                      onChange={(event) =>
                        handleFormChange("status", event.target.value)
                      }
                      disabled={isSubmitting}
                      required
                    >
                      <option value="AVAILABLE">Available</option>
                      <option value="IN_USE">In Use</option>
                      <option value="MAINTENANCE">Maintenance</option>
                    </select>
                  </div>
                )}
              </div>

              {showEditModal && selectedVehicle && !selectedVehicle.is_active && (
                <p className="request-demo-note">
                  This vehicle is inactive. Editing its information will not
                  reactivate it.
                </p>
              )}

              {formError && (
                <p className="assignment-error">{formError}</p>
              )}

              <div className="request-modal-actions">
                <button
                  type="button"
                  className="secondary-button"
                  onClick={showEditModal ? closeEditModal : closeAddModal}
                  disabled={isSubmitting}
                >
                  Cancel
                </button>

                <button
                  type="submit"
                  className="vehicles-add-button"
                  disabled={isSubmitting}
                >
                  {showEditModal ? <Pencil size={16} /> : <Plus size={17} />}
                  {isSubmitting
                    ? "Saving..."
                    : showEditModal
                      ? "Save Changes"
                      : "Save Vehicle"}
                </button>
              </div>
            </form>
          </section>
        </div>
      )}
    </div>
  );
}
