
import {
  Activity,
  Bell,
  Boxes,
  ClipboardList,
  LayoutDashboard,
  Truck,
  Users,
  Wrench,
  X,
} from "lucide-react";
import { NavLink } from "react-router-dom";
import { useAuthStore, type UserRole } from "../../stores/authStore";

interface SidebarProps {
  open: boolean;
  onClose: () => void;
}

const navigation: {
  label: string;
  items: {
    label: string;
    path: string;
    icon: typeof LayoutDashboard;
    roles: UserRole[];
  }[];
}[] = [
  {
    label: "OVERVIEW",
    items: [
      {
        label: "Dashboard",
        path: "/dashboard",
        icon: LayoutDashboard,
        roles: ["ADMIN", "DISPATCHER", "TECHNICIAN", "CUSTOMER"],
      },
    ],
  },
  {
    label: "OPERATIONS",
    items: [
      {
        label: "Service Requests",
        path: "/service-requests",
        icon: ClipboardList,
        roles: ["ADMIN", "DISPATCHER", "CUSTOMER"],
      },
      {
        label: "Work Orders",
        path: "/work-orders",
        icon: Wrench,
        roles: ["ADMIN", "DISPATCHER", "TECHNICIAN", "CUSTOMER"],
      },
      {
        label: "Technicians",
        path: "/technicians",
        icon: Users,
        roles: ["ADMIN", "DISPATCHER"],
      },
      {
        label: "Vehicles",
        path: "/vehicles",
        icon: Truck,
        roles: ["ADMIN", "DISPATCHER"],
      },
      {
        label: "Parts Inventory",
        path: "/parts",
        icon: Boxes,
        roles: ["ADMIN", "DISPATCHER"],
      },
    ],
  },
  {
    label: "SYSTEM",
    items: [
      {
        label: "Notifications",
        path: "/notifications",
        icon: Bell,
        roles: ["ADMIN", "DISPATCHER", "TECHNICIAN", "CUSTOMER"],
      },
    ],
  },
];

export default function Sidebar({ open, onClose }: SidebarProps) {
  const user = useAuthStore((state) => state.user);

  const filteredNavigation = navigation
    .map((group) => ({
      ...group,
      items: group.items.filter((item) =>
        user ? item.roles.includes(user.role) : false,
      ),
    }))
    .filter((group) => group.items.length > 0);

  return (
    <>
      {open && (
        <button
          className="sidebar-overlay"
          aria-label="Close navigation"
          onClick={onClose}
        />
      )}

      <aside className={`sidebar ${open ? "sidebar--open" : ""}`}>
        <div className="sidebar-brand">
          <div className="brand-icon">
            <Activity size={23} />
          </div>

          <div className="brand-copy">
            <span className="brand-name">SERVIXA</span>
            <span className="brand-caption">Service Management</span>
          </div>

          <button
            className="sidebar-close"
            aria-label="Close sidebar"
            onClick={onClose}
          >
            <X size={20} />
          </button>
        </div>

        <nav className="sidebar-navigation">
          {filteredNavigation.map((group) => (
            <div className="nav-group" key={group.label}>
              <p className="nav-group-label">{group.label}</p>

              {group.items.map((item) => {
                const Icon = item.icon;

                return (
                  <NavLink
                    key={item.path}
                    to={item.path}
                    onClick={onClose}
                    className={({ isActive }) =>
                      `nav-link ${isActive ? "nav-link--active" : ""}`
                    }
                  >
                    <Icon size={19} strokeWidth={1.8} />
                    <span>{item.label}</span>
                  </NavLink>
                );
              })}
            </div>
          ))}
        </nav>

        <div className="sidebar-footer">
          <div className="sidebar-footer-icon">
            <Wrench size={18} />
          </div>

          <div>
            <p className="sidebar-footer-title">SERVIXA Platform</p>
            <p className="sidebar-footer-caption">Fleet & Field Service</p>
          </div>
        </div>
      </aside>
    </>
  );
}
