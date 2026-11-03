
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, CalendarDays, CheckCircle2, ClipboardList } from "lucide-react";
import { serviceRequestApi } from "../api/serviceRequests";
import "../styles/NewServiceRequestPage.css";

export default function NewServiceRequestPage() {
  const navigate = useNavigate();

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [serviceAddress, setServiceAddress] = useState("");
  const [preferredDate, setPreferredDate] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSuccess("");
    setLoading(true);

    try {
      await serviceRequestApi.create({
        title: title.trim(),
        description: description.trim(),
        service_address: serviceAddress.trim(),
        preferred_date: preferredDate
          ? new Date(preferredDate).toISOString()
          : null,
      });

      setSuccess("Service request berhasil dibuat.");

      setTimeout(() => {
        navigate("/service-requests");
      }, 800);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Gagal membuat service request."
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="page-container new-service-request-page">
      <div className="page-heading">
        <button
          type="button"
          className="back-button"
          onClick={() => navigate("/service-requests")}
        >
          <ArrowLeft size={18} />
          Back to requests
        </button>

        <div className="section-kicker">CUSTOMER PORTAL</div>
        <h1>New Service Request</h1>
        <p>
          Describe the issue and provide your preferred service details.
        </p>
      </div>

      <form className="service-request-form" onSubmit={handleSubmit}>
        <div className="form-section-heading">
          <div className="form-section-icon">
            <ClipboardList size={20} />
          </div>
          <div>
            <h2>Request details</h2>
            <p>Tell us what needs to be serviced.</p>
          </div>
        </div>

        {error && <div className="form-alert form-alert-error">{error}</div>}
        {success && (
          <div className="form-alert form-alert-success">
            <CheckCircle2 size={18} />
            {success}
          </div>
        )}

        <div className="form-field">
          <label htmlFor="title">Request title *</label>
          <input
            id="title"
            type="text"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="e.g. Vehicle AC is not cooling"
            maxLength={150}
            required
          />
        </div>

        <div className="form-field">
          <label htmlFor="description">Problem description *</label>
          <textarea
            id="description"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            placeholder="Describe the issue in detail..."
            rows={5}
            required
          />
        </div>

        <div className="form-field">
          <label htmlFor="serviceAddress">Service address *</label>
          <input
            id="serviceAddress"
            type="text"
            value={serviceAddress}
            onChange={(event) => setServiceAddress(event.target.value)}
            placeholder="Enter the vehicle's service location"
            required
          />
        </div>

        <div className="form-field">
          <label htmlFor="preferredDate">
            <CalendarDays size={16} />
            Preferred service date
          </label>
          <input
            id="preferredDate"
            type="datetime-local"
            value={preferredDate}
            onChange={(event) => setPreferredDate(event.target.value)}
          />
          <span className="form-hint">
            Optional. Choose a date and time that works for you.
          </span>
        </div>

        <div className="form-actions">
          <button
            type="button"
            className="button-secondary"
            onClick={() => navigate("/service-requests")}
            disabled={loading}
          >
            Cancel
          </button>

          <button
            type="submit"
            className="button-primary"
            disabled={loading}
          >
            {loading ? "Submitting..." : "Submit request"}
          </button>
        </div>
      </form>
    </div>
  );
}
