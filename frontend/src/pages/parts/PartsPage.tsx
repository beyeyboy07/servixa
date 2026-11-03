
import { useMemo, useState } from "react";
import {
  AlertTriangle,
  ArrowDownToLine,
  Boxes,
  CircleDollarSign,
  ClipboardList,
  Package,
  Plus,
  Search,
  SlidersHorizontal,
  X,
} from "lucide-react";

type PartStatus = "IN_STOCK" | "LOW_STOCK" | "OUT_OF_STOCK";

type Part = {
  id: number;
  code: string;
  name: string;
  category: string;
  quantity: number;
  minimumStock: number;
  unit: string;
  unitPrice: number;
  supplier: string;
  location: string;
};

const initialParts: Part[] = [
  {
    id: 1,
    code: "PRT-001",
    name: "Engine Oil 5W-30",
    category: "Lubricants",
    quantity: 42,
    minimumStock: 10,
    unit: "Liter",
    unitPrice: 125000,
    supplier: "Auto Parts Indonesia",
    location: "Shelf A1",
  },
  {
    id: 2,
    code: "PRT-002",
    name: "Oil Filter",
    category: "Filters",
    quantity: 6,
    minimumStock: 10,
    unit: "Pcs",
    unitPrice: 85000,
    supplier: "Prima Motor Supply",
    location: "Shelf B2",
  },
  {
    id: 3,
    code: "PRT-003",
    name: "Brake Pad Set",
    category: "Braking System",
    quantity: 18,
    minimumStock: 5,
    unit: "Set",
    unitPrice: 275000,
    supplier: "Auto Parts Indonesia",
    location: "Shelf C1",
  },
  {
    id: 4,
    code: "PRT-004",
    name: "Air Filter",
    category: "Filters",
    quantity: 0,
    minimumStock: 8,
    unit: "Pcs",
    unitPrice: 95000,
    supplier: "Prima Motor Supply",
    location: "Shelf B3",
  },
  {
    id: 5,
    code: "PRT-005",
    name: "Spark Plug",
    category: "Engine",
    quantity: 24,
    minimumStock: 8,
    unit: "Pcs",
    unitPrice: 65000,
    supplier: "Bintang Automotive",
    location: "Shelf D1",
  },
  {
    id: 6,
    code: "PRT-006",
    name: "Coolant",
    category: "Fluids",
    quantity: 4,
    minimumStock: 8,
    unit: "Liter",
    unitPrice: 75000,
    supplier: "Bintang Automotive",
    location: "Shelf A2",
  },
  {
    id: 7,
    code: "PRT-007",
    name: "Drive Belt",
    category: "Engine",
    quantity: 15,
    minimumStock: 5,
    unit: "Pcs",
    unitPrice: 180000,
    supplier: "Prima Motor Supply",
    location: "Shelf D2",
  },
  {
    id: 8,
    code: "PRT-008",
    name: "Brake Fluid",
    category: "Fluids",
    quantity: 12,
    minimumStock: 5,
    unit: "Liter",
    unitPrice: 68000,
    supplier: "Auto Parts Indonesia",
    location: "Shelf A3",
  },
];

function getPartStatus(part: Part): PartStatus {
  if (part.quantity === 0) return "OUT_OF_STOCK";
  if (part.quantity <= part.minimumStock) return "LOW_STOCK";
  return "IN_STOCK";
}

const statusLabels: Record<PartStatus, string> = {
  IN_STOCK: "In Stock",
  LOW_STOCK: "Low Stock",
  OUT_OF_STOCK: "Out of Stock",
};

function formatCurrency(value: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);
}

export default function PartsPage() {
  const [parts] = useState<Part[]>(initialParts);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("ALL");
  const [categoryFilter, setCategoryFilter] = useState("ALL");
  const [selectedPart, setSelectedPart] = useState<Part | null>(null);
  const [showAddNotice, setShowAddNotice] = useState(false);

  const categories = useMemo(
    () => [...new Set(parts.map((part) => part.category))].sort(),
    [parts],
  );

  const filteredParts = useMemo(() => {
    const keyword = search.trim().toLowerCase();

    return parts.filter((part) => {
      const matchesSearch =
        !keyword ||
        part.code.toLowerCase().includes(keyword) ||
        part.name.toLowerCase().includes(keyword) ||
        part.category.toLowerCase().includes(keyword) ||
        part.supplier.toLowerCase().includes(keyword);

      const matchesStatus =
        statusFilter === "ALL" ||
        getPartStatus(part) === statusFilter;

      const matchesCategory =
        categoryFilter === "ALL" || part.category === categoryFilter;

      return matchesSearch && matchesStatus && matchesCategory;
    });
  }, [parts, search, statusFilter, categoryFilter]);

  const lowStockCount = parts.filter(
    (part) => getPartStatus(part) === "LOW_STOCK",
  ).length;

  const outOfStockCount = parts.filter(
    (part) => getPartStatus(part) === "OUT_OF_STOCK",
  ).length;

  const totalInventoryValue = parts.reduce(
    (total, part) => total + part.quantity * part.unitPrice,
    0,
  );

  function handleExport() {
    const headers = [
      "Code",
      "Part Name",
      "Category",
      "Quantity",
      "Minimum Stock",
      "Unit",
      "Unit Price",
      "Supplier",
      "Location",
      "Status",
    ];

    const rows = filteredParts.map((part) => [
      part.code,
      part.name,
      part.category,
      part.quantity,
      part.minimumStock,
      part.unit,
      part.unitPrice,
      part.supplier,
      part.location,
      statusLabels[getPartStatus(part)],
    ]);

    const csvCell = (value: string | number) =>
      `"${String(value).replace(/"/g, '""')}"`;

    const csv = [headers, ...rows]
      .map((row) => row.map(csvCell).join(","))
      .join("\r\n");

    const blob = new Blob(["\uFEFF", csv], {
      type: "text/csv;charset=utf-8;",
    });

    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");

    link.href = url;
    link.download = "servixa-parts.csv";
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
  }

  return (
    <div className="parts-page">
      <div className="page-heading">
        <div>
          <h1>Parts Inventory</h1>
          <p>Track spare parts, stock levels, and inventory value.</p>
        </div>

        <button
          className="parts-primary-button"
          onClick={() => setShowAddNotice(true)}
        >
          <Plus size={18} />
          Add Part
        </button>
      </div>

      <div className="parts-stats">
        <div className="parts-stat-card">
          <div className="parts-stat-icon blue">
            <Boxes size={21} />
          </div>
          <div>
            <span>Total Part Types</span>
            <strong>{parts.length}</strong>
            <small>Registered items</small>
          </div>
        </div>

        <div className="parts-stat-card">
          <div className="parts-stat-icon green">
            <Package size={21} />
          </div>
          <div>
            <span>Total Units</span>
            <strong>
              {parts.reduce((total, part) => total + part.quantity, 0)}
            </strong>
            <small>Across all inventory</small>
          </div>
        </div>

        <div className="parts-stat-card">
          <div className="parts-stat-icon orange">
            <AlertTriangle size={21} />
          </div>
          <div>
            <span>Low Stock</span>
            <strong>{lowStockCount}</strong>
            <small>Restocking recommended</small>
          </div>
        </div>

        <div className="parts-stat-card">
          <div className="parts-stat-icon purple">
            <CircleDollarSign size={21} />
          </div>
          <div>
            <span>Inventory Value</span>
            <strong className="parts-value-number">
              {formatCurrency(totalInventoryValue)}
            </strong>
            <small>Based on current stock</small>
          </div>
        </div>
      </div>

      {outOfStockCount > 0 && (
        <div className="parts-alert">
          <AlertTriangle size={19} />
          <div>
            <strong>
              {outOfStockCount} item{outOfStockCount > 1 ? "s" : ""} out of
              stock
            </strong>
            <p>Review unavailable items before assigning parts to work orders.</p>
          </div>
        </div>
      )}

      <section className="parts-panel">
        <div className="parts-toolbar">
          <div>
            <h2>Inventory List</h2>
            <p>Search and filter your spare parts inventory.</p>
          </div>

          <button
            className="parts-export-button"
            onClick={handleExport}
          >
            <ArrowDownToLine size={16} />
            Export CSV
          </button>
        </div>

        <div className="parts-filters">
          <label className="work-orders-search">
            <Search size={17} />
            <input
              type="search"
              placeholder="Search parts, code, supplier..."
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </label>

          <label className="parts-filter-select">
            <SlidersHorizontal size={16} />
            <select
              aria-label="Filter by stock status"
              value={statusFilter}
              onChange={(event) => setStatusFilter(event.target.value)}
            >
              <option value="ALL">All Stock Status</option>
              <option value="IN_STOCK">In Stock</option>
              <option value="LOW_STOCK">Low Stock</option>
              <option value="OUT_OF_STOCK">Out of Stock</option>
            </select>
          </label>

          <label className="parts-filter-select">
            <ClipboardList size={16} />
            <select
              aria-label="Filter by category"
              value={categoryFilter}
              onChange={(event) => setCategoryFilter(event.target.value)}
            >
              <option value="ALL">All Categories</option>
              {categories.map((category) => (
                <option value={category} key={category}>
                  {category}
                </option>
              ))}
            </select>
          </label>
        </div>

        <div className="parts-table-wrapper">
          <table className="parts-table">
            <thead>
              <tr>
                <th>Part</th>
                <th>Category</th>
                <th>Stock Level</th>
                <th>Unit Price</th>
                <th>Supplier</th>
                <th>Status</th>
                <th>Action</th>
              </tr>
            </thead>

            <tbody>
              {filteredParts.map((part) => {
                const status = getPartStatus(part);
                const stockPercent = Math.min(
                  100,
                  (part.quantity / Math.max(part.minimumStock * 2, 1)) * 100,
                );

                return (
                  <tr key={part.id}>
                    <td>
                      <div className="parts-item-cell">
                        <div className="parts-item-icon">
                          <Package size={18} />
                        </div>
                        <div>
                          <strong>{part.name}</strong>
                          <span>{part.code}</span>
                        </div>
                      </div>
                    </td>

                    <td>
                      <span className="parts-category">{part.category}</span>
                    </td>

                    <td>
                      <div className="parts-stock-cell">
                        <div>
                          <strong>{part.quantity}</strong>
                          <span>
                            {part.unit} · Min. {part.minimumStock}
                          </span>
                        </div>
                        <div className="parts-stock-track">
                          <span
                            className={`parts-stock-fill ${status.toLowerCase()}`}
                            style={{ width: `${stockPercent}%` }}
                          />
                        </div>
                      </div>
                    </td>

                    <td className="parts-price">
                      {formatCurrency(part.unitPrice)}
                    </td>

                    <td>
                      <span className="parts-supplier">{part.supplier}</span>
                    </td>

                    <td>
                      <span className={`parts-status ${status.toLowerCase()}`}>
                        {statusLabels[status]}
                      </span>
                    </td>

                    <td>
                      <button
                        className="parts-view-button"
                        onClick={() => setSelectedPart(part)}
                      >
                        View
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        {filteredParts.length === 0 && (
          <div className="parts-empty">
            <Package size={32} />
            <h3>No parts found</h3>
            <p>Try changing your search or filters.</p>
          </div>
        )}

        <div className="parts-results">
          Showing {filteredParts.length} of {parts.length} part types
        </div>
      </section>

      {selectedPart && (
        <div
          className="request-modal-backdrop"
          onClick={() => setSelectedPart(null)}
        >
          <section
            className="request-modal parts-detail-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="part-detail-title"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="request-modal-heading">
              <div>
                <h2 id="part-detail-title">{selectedPart.name}</h2>
                <p>{selectedPart.code}</p>
              </div>
              <button
                className="parts-close-button"
                onClick={() => setSelectedPart(null)}
                aria-label="Close part details"
              >
                <X size={20} />
              </button>
            </div>

            <div className="parts-detail-summary">
              <div className="parts-detail-icon">
                <Package size={28} />
              </div>
              <div>
                <span>Current stock</span>
                <strong>
                  {selectedPart.quantity} {selectedPart.unit}
                </strong>
                <span
                  className={`parts-status ${getPartStatus(selectedPart).toLowerCase()}`}
                >
                  {statusLabels[getPartStatus(selectedPart)]}
                </span>
              </div>
            </div>

            <div className="request-detail-grid">
              <div>
                <span>Category</span>
                <strong>{selectedPart.category}</strong>
              </div>
              <div>
                <span>Minimum Stock</span>
                <strong>
                  {selectedPart.minimumStock} {selectedPart.unit}
                </strong>
              </div>
              <div>
                <span>Unit Price</span>
                <strong>{formatCurrency(selectedPart.unitPrice)}</strong>
              </div>
              <div>
                <span>Stock Value</span>
                <strong>
                  {formatCurrency(
                    selectedPart.quantity * selectedPart.unitPrice,
                  )}
                </strong>
              </div>
              <div>
                <span>Supplier</span>
                <strong>{selectedPart.supplier}</strong>
              </div>
              <div>
                <span>Storage Location</span>
                <strong>{selectedPart.location}</strong>
              </div>
            </div>

            <div className="request-demo-note">
              Inventory information currently uses sample data. Stock
              adjustments and editing will be implemented in a later step.
            </div>

            <div className="request-modal-actions">
              <button
                className="secondary-button"
                onClick={() => setSelectedPart(null)}
              >
                Close
              </button>
            </div>
          </section>
        </div>
      )}

      {showAddNotice && (
        <div
          className="request-modal-backdrop"
          onClick={() => setShowAddNotice(false)}
        >
          <section
            className="request-modal"
            role="dialog"
            aria-modal="true"
            aria-labelledby="part-add-title"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="request-modal-heading">
              <div>
                <h2 id="part-add-title">Add Part</h2>
                <p>Register a new spare part in your inventory.</p>
              </div>
              <button
                className="parts-close-button"
                onClick={() => setShowAddNotice(false)}
                aria-label="Close"
              >
                <X size={20} />
              </button>
            </div>

            <div className="request-demo-note">
              The part registration form and backend integration will be added
              in a later step.
            </div>

            <div className="request-modal-actions">
              <button
                className="secondary-button"
                onClick={() => setShowAddNotice(false)}
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
