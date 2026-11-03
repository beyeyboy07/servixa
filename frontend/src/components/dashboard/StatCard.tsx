import type { LucideIcon } from "lucide-react";
import type { ReactNode } from "react";

interface StatCardProps {
  title: string;
  value: string;
  description: ReactNode;
  icon: LucideIcon;
  tone: "blue" | "green" | "orange" | "purple";
}

export default function StatCard({
  title,
  value,
  description,
  icon: Icon,
  tone,
}: StatCardProps) {
  return (
    <article className="stat-card">
      <div className={`stat-card-icon stat-card-icon--${tone}`}>
        <Icon size={21} strokeWidth={1.9} />
      </div>

      <div className="stat-card-content">
        <p className="stat-card-title">{title}</p>
        <p className="stat-card-value">{value}</p>
        <p className="stat-card-description">{description}</p>
      </div>
    </article>
  );
}