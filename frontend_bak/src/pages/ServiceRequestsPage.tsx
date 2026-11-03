
import { useEffect, useState } from "react";
import { ArrowRight, ClipboardList, RefreshCw } from "lucide-react";

import {
  serviceRequestApi,
  type ServiceRequest,
} from "../api/serviceRequests";

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

function getStatusClass(status: string) {
  return `status-${status.toLowerCase().replaceAll("_", "-")}`;
}

export default function ServiceRequestsPage() {
  const [requests, setRequests] = useState<ServiceRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function loadRequests() {
    setLoading(true);
    setError("");

    try {
      const data = await serviceRequestApi.getAll();
      setRequests(data);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to load service requests.",
      );
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadRequests();
  }, []);

  return (
    <div className="dashboard">
      <div className="page-heading">
        <div>
          <div className="page-eyebrow">SERVICE OPERATIONS</div>
          <h1>Service Requests</h1>
          <p>Review and manage customer service requests.</p>
        </div>

        <button
          className="text-button"
          onClick={() => void loadRequests()}
          disabled={loading}
        >
          <RefreshCw size={16} />
          Refresh
        </button>
      </div>

      <section className="stats-grid">
        <article className="stat-card">
          <div className="stat-card-top">
            <span>Total Requests</span>
            <div className="stat-icon stat-icon-blue">
              <ClipboardList size={20} />
            </div>
          </div>

          <div className="stat-value">
            {loading ? "—" : requests.length.toLocaleString()}
          </div>

          <div className="stat-foot">
            <span>Requests returned by API</span>
          </div>
        </article>
      </section>

      <section className="panel table-panel">
        <div className="panel-heading">
          <div>
            <h2>All Service Requests</h2>
            <p>Customer requests submitted to SERVIXA</p>
          </div>
        </div>

        {loading ? (
          <p className="dashboard-table-message">
            Loading service requests...
          </p>
        ) : error ? (
          <div className="dashboard-error" role="alert">
            <p>Unable to load service requests: {error}</p>
            <button onClick={() => void loadRequests()}>
              Retry
            </button>
          </div>
        ) : requests.length === 0 ? (
          <p className="dashboard-table-message">
            No service requests found.
          </p>
        ) : (
          <div className="table-scroll">
            <table className="work-orders-table">
              
            <thead>
                <tr>
                    <th>Request</th>
                    <th>Customer ID</th>
                    <th>Service Address</th>
                    <th>Preferred Date</th>
                    <th>Created</th>
                    <th>Status</th>
                    <th />
                </tr>
            </thead>


              <tbody>
                {requests.map((request) => (
                  <tr key={request.id}>
                    <td>
                      <div className="order-id">
                        SR-{request.id}
                      </div>
                      <div className="order-title">
                        {request.title}
                      </div>
                    </td>

                    <td>#{request.customer_id}</td>

                    <td>{request.service_address || "—"}</td>
                    
                    <td>
                        {request.preferred_date
                        ? formatDate(request.preferred_date)
                        : "—"}
                    </td>


                    <td>{formatDate(request.created_at)}</td>

                    <td>
                      <span
                        className={`status-badge ${getStatusClass(
                          request.status,
                        )}`}
                      >
                        <i />
                        {formatStatus(request.status)}
                      </span>
                    </td>

                    <td>
                      <button
                        className="row-action"
                        aria-label={`Open service request ${request.id}`}
                        onClick={() =>
                          window.location.assign(
                            `/service-requests/${request.id}`,
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
    </div>
  );
}
