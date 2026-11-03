
import { Wrench } from "lucide-react";

interface FeaturePlaceholderPageProps {
  title: string;
  description?: string;
}

export default function FeaturePlaceholderPage({
  title,
  description = "This feature will be connected to the SERVIXA API in the next implementation step.",
}: FeaturePlaceholderPageProps) {
  return (
    <div className="dashboard">
      <div className="page-heading">
        <div>
          <div className="page-eyebrow">SERVIXA WORKSPACE</div>
          <h1>{title}</h1>
          <p>{description}</p>
        </div>
      </div>

      <section className="panel" style={{ padding: "32px" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
          <div className="stat-icon stat-icon-blue">
            <Wrench size={22} />
          </div>

          <div>
            <h2 style={{ margin: 0 }}>Coming soon</h2>
            <p style={{ marginTop: "8px" }}>
              Halaman ini sudah terdaftar berdasarkan role akunmu.
              Fungsionalitas dan integrasi API akan ditambahkan
              pada tahap berikutnya.
            </p>
          </div>
        </div>
      </section>
    </div>
  );
}
