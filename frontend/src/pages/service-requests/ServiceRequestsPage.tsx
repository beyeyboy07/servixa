
import { useEffect, useMemo, useState } from "react";
import {
  AlertCircle,
  CalendarDays,
  CheckCircle2,
  ClipboardList,
  Clock3,
  Eye,
  Plus,
  Search,
  SlidersHorizontal,
  XCircle,
} from "lucide-react";
import Swal from "sweetalert2";
import { useAuthStore } from "../../stores/authStore";
import {
  createServiceRequest,
  getMyServiceRequests,
  type ServiceRequest,
  type ServiceRequestStatus,
} from "../../services/serviceRequestService";

const statusLabels: Record<ServiceRequestStatus, string> = {
  PENDING: "Pending",
  REVIEWED: "Reviewed",
  REJECTED: "Rejected",
};

function RequestStatusBadge({
  status,
}: {
  status: ServiceRequestStatus;
}) {
  return (
    <span
      className={`request-status-badge request-status-badge--${status.toLowerCase().replaceAll("_", "-")}`}
    >
      <span className="status-dot" />
      {statusLabels[status]}
    </span>
  );
}

function formatRequestId(id: number) {
  return `SR-${String(id).padStart(5, "0")}`;
}

function formatDate(date: string | null) {
  if (!date) return "-";

  return new Date(date).toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

function getErrorMessage(error: unknown) {
  if (error instanceof Error) return error.message;
  return "Something went wrong. Please try again.";
}

export default function ServiceRequestsPage() {
  const { user } = useAuthStore();

  const isCustomer = user?.role === "CUSTOMER";
  const canManage =
    user?.role === "ADMIN" || user?.role === "DISPATCHER";

  const [requests, setRequests] = useState<ServiceRequest[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("ALL");

  const [selectedRequest, setSelectedRequest] =
    useState<ServiceRequest | null>(null);

  const [showCreateModal, setShowCreateModal] = useState(false);

  const [form, setForm] = useState({
    title: "",
    description: "",
    service_address: "",
    preferred_date: "",
  });

  const loadRequests = async () => {
    if (!isCustomer) {
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError("");

      const data = await getMyServiceRequests();
      setRequests(data);
    } catch (err) {
      setError(getErrorMessage(err));
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadRequests();
  }, [user?.role]);

  const handleCreateRequest = async (
    event: React.FormEvent<HTMLFormElement>,
  ) => {
    event.preventDefault();

    try {
      setIsSubmitting(true);

      const payload = {
        title: form.title.trim(),
        description: form.description.trim(),
        service_address: form.service_address.trim(),
        preferred_date: form.preferred_date
          ? `${form.preferred_date}T00:00:00+07:00`
          : null,
      };

      const createdRequest = await createServiceRequest(payload);

      setRequests((previous) => [createdRequest, ...previous]);

      setShowCreateModal(false);
      setForm({
        title: "",
        description: "",
        service_address: "",
        preferred_date: "",
      });

      await Swal.fire({
        icon: "success",
        title: "Request Submitted",
        text: "Your service request has been submitted successfully.",
        confirmButtonText: "OK",
      });
    } catch (err) {
      Swal.fire({
        icon: "error",
        title: "Failed to Submit",
        text: getErrorMessage(err),
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  const filteredRequests = useMemo(() => {
    const keyword = search.trim().toLowerCase();

    return requests.filter((request) => {
      const matchesSearch =
        !keyword ||
        formatRequestId(request.id).toLowerCase().includes(keyword) ||
        request.title.toLowerCase().includes(keyword) ||
        request.description.toLowerCase().includes(keyword) ||
        request.service_address.toLowerCase().includes(keyword);

      const matchesStatus =
        statusFilter === "ALL" || request.status === statusFilter;

      return matchesSearch && matchesStatus;
    });
  }, [requests, search, statusFilter]);

  const counts = {
    total: requests.length,
    pending: requests.filter((item) => item.status === "PENDING").length,
    reviewed: requests.filter((item) => item.status === "REVIEWED").length,
    rejected: requests.filter((item) => item.status === "REJECTED").length,
  };

  return (
    <div className="service-requests-page">
      <section className="service-requests-heading">
        <div>
          <p className="dashboard-eyebrow">
            OPERATIONS / SERVICE REQUESTS
          </p>
          <h1 className="page-title">Service Requests</h1>
          <p className="page-description">
            {isCustomer
              ? "Submit service requests and track their current status."
              : "Review and manage customer service requests."}
          </p>
        </div>

        {isCustomer && (
          <button
            type="button"
            className="quick-action-button quick-action-button--primary"
            onClick={() => setShowCreateModal(true)}
          >
            <Plus size={17} />
            New Service Request
          </button>
        )}
      </section>

      {isCustomer && (
        <section className="service-request-stats">
          <article className="service-request-stat">
            <div className="request-stat-icon request-stat-icon--blue">
              <ClipboardList size={20} />
            </div>
            <div>
              <p>Total Requests</p>
              <strong>{counts.total}</strong>
            </div>
          </article>

          <article className="service-request-stat">
            <div className="request-stat-icon request-stat-icon--orange">
              <Clock3 size={20} />
            </div>
            <div>
              <p>Pending</p>
              <strong>{counts.pending}</strong>
            </div>
          </article>

          <article className="service-request-stat">
            <div className="request-stat-icon request-stat-icon--green">
              <CheckCircle2 size={20} />
            </div>
            <div>
              <p>Reviewed</p>
              <strong>{counts.reviewed}</strong>
            </div>
          </article>

          <article className="service-request-stat">
            <div className="request-stat-icon request-stat-icon--purple">
              <XCircle size={20} />
            </div>
            <div>
              <p>Rejected</p>
              <strong>{counts.rejected}</strong>
            </div>
          </article>
        </section>
      )}

      <section className="dashboard-panel service-request-list-panel">
        <div className="service-request-list-heading">
          <div>
            <h2 className="panel-title">My Service Requests</h2>
            <p className="panel-description">
              {filteredRequests.length} of {requests.length} requests
            </p>
          </div>
        </div>

        <div className="service-request-toolbar">
          <label className="work-orders-search">
            <Search size={18} />
            <input
              type="search"
              placeholder="Search request, title, description..."
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </label>

          <label className="work-orders-filter">
            <SlidersHorizontal size={17} />
            <select
              value={statusFilter}
              onChange={(event) => setStatusFilter(event.target.value)}
              aria-label="Filter service requests by status"
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

        {isLoading ? (
          <div className="work-orders-empty">
            <Clock3 size={26} />
            <strong>Loading service requests...</strong>
          </div>
        ) : error ? (
          <div className="work-orders-empty">
            <AlertCircle size={26} />
            <strong>Failed to load service requests</strong>
            <p>{error}</p>
            <button
              type="button"
              className="text-action"
              onClick={loadRequests}
            >
              Try again
            </button>
          </div>
        ) : !isCustomer ? (
          <div className="work-orders-empty">
            <AlertCircle size={26} />
            <strong>Management list endpoint is not available yet</strong>
            <p>
              The backend currently provides the request list endpoint
              only for customers. ADMIN and DISPATCHER need a management
              GET endpoint before they can load all customer requests here.
            </p>
          </div>
        ) : (
          <>
            <div className="table-container">
              <table className="work-orders-table service-request-table">
                <thead>
                  <tr>
                    <th>Request</th>
                    <th>Service</th>
                    <th>Service Address</th>
                    <th>Preferred Date</th>
                    <th>Submitted</th>
                    <th>Status</th>
                    <th>Action</th>
                  </tr>
                </thead>

                <tbody>
                  {filteredRequests.map((request) => (
                    <tr key={request.id}>
                      <td>
                        <span className="work-order-id">
                          {formatRequestId(request.id)}
                        </span>
                      </td>

                      <td>
                        <div className="table-primary-text">
                          {request.title}
                        </div>
                        <div className="table-secondary-text">
                          {request.description}
                        </div>
                      </td>

                      <td>{request.service_address}</td>

                      <td>{formatDate(request.preferred_date)}</td>

                      <td>
                        <span className="work-order-date">
                          <CalendarDays size={14} />
                          {formatDate(request.created_at)}
                        </span>
                      </td>

                      <td>
                        <RequestStatusBadge status={request.status} />
                      </td>

                      <td>
                        <button
                          className="request-view-button"
                          type="button"
                          onClick={() => setSelectedRequest(request)}
                        >
                          <Eye size={15} />
                          View
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>

              {filteredRequests.length === 0 && (
                <div className="work-orders-empty">
                  <Search size={26} />
                  <strong>No service requests found</strong>
                  <p>
                    {requests.length === 0
                      ? "You haven't submitted any service requests yet."
                      : "Try another search keyword or status filter."}
                  </p>

                  {requests.length > 0 && (
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
              <span>{filteredRequests.length} results</span>
            </div>
          </>
        )}
      </section>

      {showCreateModal && (
        <div
          className="request-modal-backdrop"
          onClick={() => {
            if (!isSubmitting) setShowCreateModal(false);
          }}
        >
          <section
            className="request-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="create-request-title"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="request-modal-heading">
              <div>
                <p className="dashboard-eyebrow">
                  NEW SERVICE REQUEST
                </p>
                <h2 id="create-request-title">Request a Service</h2>
              </div>

              <button
                className="panel-more-button"
                type="button"
                aria-label="Close form"
                disabled={isSubmitting}
                onClick={() => setShowCreateModal(false)}
              >
                <XCircle size={20} />
              </button>
            </div>

            <form onSubmit={handleCreateRequest}>
              <div className="request-form-fields">
                <label>
                  <span>Service Title *</span>
                  <input
                    type="text"
                    value={form.title}
                    onChange={(event) =>
                      setForm({
                        ...form,
                        title: event.target.value,
                      })
                    }
                    placeholder="e.g. Engine Maintenance"
                    maxLength={150}
                    required
                  />
                </label>

                <label>
                  <span>Description *</span>
                  <textarea
                    value={form.description}
                    onChange={(event) =>
                      setForm({
                        ...form,
                        description: event.target.value,
                      })
                    }
                    placeholder="Describe the problem or service you need..."
                    rows={4}
                    required
                  />
                </label>

                <label>
                  <span>Service Address *</span>
                  <textarea
                    value={form.service_address}
                    onChange={(event) =>
                      setForm({
                        ...form,
                        service_address: event.target.value,
                      })
                    }
                    placeholder="Enter the location where service is required..."
                    rows={3}
                    required
                  />
                </label>

                <label>
                  <span>Preferred Date</span>
                  <input
                    type="date"
                    value={form.preferred_date}
                    onChange={(event) =>
                      setForm({
                        ...form,
                        preferred_date: event.target.value,
                      })
                    }
                  />
                </label>
              </div>

              <div className="request-modal-actions">
                <button
                  className="quick-action-button"
                  type="button"
                  disabled={isSubmitting}
                  onClick={() => setShowCreateModal(false)}
                >
                  Cancel
                </button>

                <button
                  className="quick-action-button quick-action-button--primary"
                  type="submit"
                  disabled={isSubmitting}
                >
                  <Plus size={16} />
                  {isSubmitting ? "Submitting..." : "Submit Request"}
                </button>
              </div>
            </form>
          </section>
        </div>
      )}

      {selectedRequest && (
        <div
          className="request-modal-backdrop"
          onClick={() => setSelectedRequest(null)}
        >
          <section
            className="request-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="request-modal-title"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="request-modal-heading">
              <div>
                <p className="dashboard-eyebrow">
                  SERVICE REQUEST DETAILS
                </p>
                <h2 id="request-modal-title">
                  {formatRequestId(selectedRequest.id)}
                </h2>
              </div>

              <button
                className="panel-more-button"
                type="button"
                aria-label="Close details"
                onClick={() => setSelectedRequest(null)}
              >
                <XCircle size={20} />
              </button>
            </div>

            <div className="request-modal-status">
              <RequestStatusBadge status={selectedRequest.status} />
            </div>

            <div className="request-detail-grid">
              <div>
                <span>Service Title</span>
                <strong>{selectedRequest.title}</strong>
              </div>

              <div>
                <span>Submitted Date</span>
                <strong>{formatDate(selectedRequest.created_at)}</strong>
              </div>

              <div>
                <span>Preferred Date</span>
                <strong>{formatDate(selectedRequest.preferred_date)}</strong>
              </div>

              <div>
                <span>Service Address</span>
                <strong>{selectedRequest.service_address}</strong>
              </div>
            </div>

            <div className="request-description-box">
              <span>Request Description</span>
              <p>{selectedRequest.description}</p>
            </div>

            <div className="request-modal-actions">
              <button
                className="quick-action-button"
                type="button"
                onClick={() => setSelectedRequest(null)}
              >
                Close
              </button>
            </div>
          </section>
        </div>
      )}
    </div>
  );
}
