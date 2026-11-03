import type { LucideIcon } from "lucide-react";
import { ArrowUpRight, Construction } from "lucide-react";

type PlaceholderPageProps = {
  title: string;
  description: string;
  icon: LucideIcon;
};

export default function PlaceholderPage({
  title,
  description,
  icon: Icon,
}: PlaceholderPageProps) {
  return (
    <div className="placeholder-page">
      <div className="page-heading">
        <div>
          <div className="eyebrow">
            <Construction size={14} />
            SERVIXA MODULE
          </div>

          <h1>{title}</h1>
          <p className="page-description">{description}</p>
        </div>
      </div>

      <section className="placeholder-panel">
        <div className="placeholder-illustration">
          <Icon size={34} strokeWidth={1.6} />
        </div>

        <span className="placeholder-badge">IN DEVELOPMENT</span>

        <h2>{title} workspace</h2>

        <p>
          This module is part of the SERVIXA operations platform.
          Its interface and functionality will be implemented in
          the next development stages.
        </p>

        <div className="placeholder-divider" />

        <div className="placeholder-info">
          <span>Module status</span>
          <strong>
            <span className="status-dot" />
            UI foundation
          </strong>
        </div>

        <div className="placeholder-info">
          <span>Backend integration</span>
          <strong>Pending</strong>
        </div>

        <div className="placeholder-info">
          <span>Data</span>
          <strong>Not connected</strong>
        </div>

        <div className="placeholder-footer">
          <span>Part of SERVIXA Fleet & Field Service Management</span>
          <ArrowUpRight size={16} />
        </div>
      </section>
    </div>
  );
}