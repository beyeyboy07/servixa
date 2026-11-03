
import { useMemo, useState } from "react";
import {
  CheckCircle2,
  Clock3,
  Search,
  SlidersHorizontal,
  UserRound,
  Users,
  Wrench,
  X,
} from "lucide-react";

type TechnicianStatus = "AVAILABLE" | "ON_JOB" | "OFF_DUTY";

interface Technician {
  id: string;
  name: string;
  email: string;
  phone: string;
  specialization: string;
  status: TechnicianStatus;
  completedJobs: number;
  activeJobs: number;
}

const initialTechnicians: Technician[] = [
  {
    id: "TECH-001",
    name: "Andi Pratama",
    email: "andi@example.com",
    phone: "0812-0000-1001",
    specialization: "Engine Maintenance",
    status: "ON_JOB",
    completedJobs: 42,
    activeJobs: 2,
  },
  {
    id: "TECH-002",
    name: "Rizky Saputra",
    email: "rizky@example.com",
    phone: "0812-0000-1002",
    specialization: "Brake & Suspension",
    status: "AVAILABLE",
    completedJobs: 37,
    activeJobs: 0,
  },
  {
    id: "TECH-003",
    name: "Dimas Setiawan",
    email: "dimas@example.com",
    phone: "0812-0000-1003",
    specialization: "General Maintenance",
    status: "ON_JOB",
    completedJobs: 31,
    activeJobs: 1,
  },
  {
    id: "TECH-004",
    name: "Budi Santoso",
    email: "budi@example.com",
    phone: "0812-0000-1004",
    specialization: "Vehicle Diagnostics",
    status: "AVAILABLE",
    completedJobs: 29,
    activeJobs: 0,
  },
  {
    id: "TECH-005",
    name: "Fajar Hidayat",
    email: "fajar@example.com",
    phone: "0812-0000-1005",
    specialization: "Tire & Wheel Service",
    status: "OFF_DUTY",
    completedJobs: 24,
    activeJobs: 0,
  },
  {
    id: "TECH-006",
    name: "Agus Firmansyah",
    email: "agus@example.com",
    phone: "0812-0000-1006",
    specialization: "Electrical Systems",
    status: "AVAILABLE",
    completedJobs: 18,
    activeJobs: 0,
  },
];

const statusLabels: Record<TechnicianStatus, string> = {
  AVAILABLE: "Available",
  ON_JOB: "On Job",
  OFF_DUTY: "Off Duty",
};

function TechnicianStatusBadge({ status }: { status: TechnicianStatus }) {
  return (
    <span
      className={`technician-status technician-status--${status.toLowerCase().replaceAll("_", "-")}`}
    >
      <span className="status-dot" />
      {statusLabels[status]}
    </span>
  );
}

function getInitials(name: string) {
  return name
    .split(" ")
    .map((part) => part[0])
    .slice(0, 2)
    .join("");
}

export default function TechniciansPage() {
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("ALL");
  const [selectedTechnician, setSelectedTechnician] =
    useState<Technician | null>(null);

  const technicians = useMemo(() => {
    const keyword = search.trim().toLowerCase();

    return initialTechnicians.filter((technician) => {
      const matchesSearch =
        !keyword ||
        technician.name.toLowerCase().includes(keyword) ||
        technician.email.toLowerCase().includes(keyword) ||
        technician.specialization.toLowerCase().includes(keyword) ||
        technician.id.toLowerCase().includes(keyword);

      const matchesStatus =
        statusFilter === "ALL" || technician.status === statusFilter;

      return matchesSearch && matchesStatus;
    });
  }, [search, statusFilter]);

  const availableCount = initialTechnicians.filter(
    (item) => item.status === "AVAILABLE",
  ).length;

  const onJobCount = initialTechnicians.filter(
    (item) => item.status === "ON_JOB",
  ).length;

  const activeJobCount = initialTechnicians.reduce(
    (total, item) => total + item.activeJobs,
    0,
  );

  return (
    <div className="technicians-page">
      <section className="technicians-heading">
        <div>
          <p className="dashboard-eyebrow">OPERATIONS / TECHNICIANS</p>
          <h1 className="page-title">Technicians</h1>
          <p className="page-description">
            Manage your field team and monitor technician workloads.
          </p>
        </div>
      </section>

      <section className="technician-stats">
        <article className="technician-stat-card">
          <div className="technician-stat-icon technician-stat-icon--blue">
            <Users size={21} />
          </div>
          <div>
            <p>Total Technicians</p>
            <strong>{initialTechnicians.length}</strong>
          </div>
        </article>

        <article className="technician-stat-card">
          <div className="technician-stat-icon technician-stat-icon--green">
            <CheckCircle2 size={21} />
          </div>
          <div>
            <p>Available</p>
            <strong>{availableCount}</strong>
          </div>
        </article>

        <article className="technician-stat-card">
          <div className="technician-stat-icon technician-stat-icon--orange">
            <Wrench size={21} />
          </div>
          <div>
            <p>Currently on Job</p>
            <strong>{onJobCount}</strong>
          </div>
        </article>

        <article className="technician-stat-card">
          <div className="technician-stat-icon technician-stat-icon--purple">
            <Clock3 size={21} />
          </div>
          <div>
            <p>Active Assignments</p>
            <strong>{activeJobCount}</strong>
          </div>
        </article>
      </section>

      <section className="dashboard-panel technicians-list-panel">
        <div className="technicians-list-heading">
          <div>
            <h2 className="panel-title">Technician Directory</h2>
            <p className="panel-description">
              {technicians.length} of {initialTechnicians.length} sample
              technicians
            </p>
          </div>
        </div>

        <div className="technicians-toolbar">
          <label className="work-orders-search">
            <Search size={18} />
            <input
              type="search"
              placeholder="Search name, email, skill, or ID..."
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </label>

          <label className="work-orders-filter">
            <SlidersHorizontal size={17} />
            <select
              aria-label="Filter technicians by status"
              value={statusFilter}
              onChange={(event) => setStatusFilter(event.target.value)}
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

        <div className="technician-grid">
          {technicians.map((technician) => (
            <article className="technician-card" key={technician.id}>
              <div className="technician-card-top">
                <div className="technician-profile">
                  <div className="technician-large-avatar">
                    {getInitials(technician.name)}
                  </div>
                  <div className="technician-profile-info">
                    <h3>{technician.name}</h3>
                    <p>{technician.id}</p>
                  </div>
                </div>
                <TechnicianStatusBadge status={technician.status} />
              </div>

              <div className="technician-specialization">
                <Wrench size={15} />
                {technician.specialization}
              </div>

              <div className="technician-contact">
                <span>{technician.email}</span>
                <span>{technician.phone}</span>
              </div>

              <div className="technician-workload">
                <div>
                  <span>Completed jobs</span>
                  <strong>{technician.completedJobs}</strong>
                </div>
                <div>
                  <span>Active jobs</span>
                  <strong>{technician.activeJobs}</strong>
                </div>
              </div>

              <button
                className="technician-detail-button"
                type="button"
                onClick={() => setSelectedTechnician(technician)}
              >
                <UserRound size={16} />
                View Profile
              </button>
            </article>
          ))}
        </div>

        {technicians.length === 0 && (
          <div className="work-orders-empty">
            <Search size={26} />
            <strong>No technicians found</strong>
            <p>Try another search keyword or status filter.</p>
            <button
              className="text-action"
              type="button"
              onClick={() => {
                setSearch("");
                setStatusFilter("ALL");
              }}
            >
              Clear filters
            </button>
          </div>
        )}

        <div className="work-orders-table-footer">
          <span>{technicians.length} results</span>
          <span className="work-orders-demo-label">Demo data</span>
        </div>
      </section>

      {selectedTechnician && (
        <div
          className="request-modal-backdrop"
          onClick={() => setSelectedTechnician(null)}
        >
          <section
            className="request-modal technician-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="technician-modal-title"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="request-modal-heading">
              <div>
                <p className="dashboard-eyebrow">TECHNICIAN PROFILE</p>
                <h2 id="technician-modal-title">
                  {selectedTechnician.name}
                </h2>
              </div>
              <button
                className="panel-more-button"
                type="button"
                aria-label="Close technician profile"
                onClick={() => setSelectedTechnician(null)}
              >
                <X size={20} />
              </button>
            </div>

            <div className="technician-modal-profile">
              <div className="technician-large-avatar technician-modal-avatar">
                {getInitials(selectedTechnician.name)}
              </div>
              <div>
                <TechnicianStatusBadge status={selectedTechnician.status} />
                <p>{selectedTechnician.id}</p>
              </div>
            </div>

            <div className="request-detail-grid">
              <div>
                <span>Email</span>
                <strong>{selectedTechnician.email}</strong>
              </div>
              <div>
                <span>Phone</span>
                <strong>{selectedTechnician.phone}</strong>
              </div>
              <div>
                <span>Specialization</span>
                <strong>{selectedTechnician.specialization}</strong>
              </div>
              <div>
                <span>Completed jobs</span>
                <strong>{selectedTechnician.completedJobs}</strong>
              </div>
              <div>
                <span>Active assignments</span>
                <strong>{selectedTechnician.activeJobs}</strong>
              </div>
            </div>

            <p className="request-demo-note">
              Profile information is sample data and will be replaced with
              technician records from the backend.
            </p>

            <div className="request-modal-actions">
              <button
                className="quick-action-button"
                type="button"
                onClick={() => setSelectedTechnician(null)}
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
