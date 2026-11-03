
import {
  LayoutDashboard,
  ClipboardList,
  Wrench,
  Users,
  Truck,
  Package,
  Bell,
  Settings,
  ChevronLeft,
  CalendarDays,
  History,
  Search,
  PlusCircle,
} from "lucide-react";
import { NavLink } from "react-router-dom";
import { useState } from "react";
import { useAuthStore, type AuthUser } from "../stores/authStore";

type UserRole = AuthUser["role"];

interface NavigationItem {
  label: string;
  path: string;
  icon: typeof LayoutDashboard;
}

interface NavigationGroup {
  label: string;
  items: NavigationItem[];
}

const dashboardItem: NavigationItem = {
  label: "Dashboard",
  path: "/dashboard",
  icon: LayoutDashboard,
};

const notificationItem: NavigationItem = {
  label: "Notifications",
  path: "/notifications",
  icon: Bell,
};

const settingsItem: NavigationItem = {
  label: "Settings",
  path: "/settings",
  icon: Settings,
};

const navigationByRole: Record<UserRole, NavigationGroup[]> = {
  ADMIN: [
    {
      label: "WORKSPACE",
      items: [
        dashboardItem,
        {
          label: "Work Orders",
          path: "/work-orders",
          icon: Wrench,
        },
      ],
    },
    {
      label: "OPERATIONS",
      items: [
        {
          label: "Technicians",
          path: "/technicians",
          icon: Users,
        },
        {
          label: "Vehicles",
          path: "/vehicles",
          icon: Truck,
        },
        {
          label: "Parts & Inventory",
          path: "/parts",
          icon: Package,
        },
      ],
    },
    {
      label: "GENERAL",
      items: [notificationItem, settingsItem],
    },
  ],

  DISPATCHER: [
    {
      label: "WORKSPACE",
      items: [
        dashboardItem,
        {
          label: "Work Orders",
          path: "/work-orders",
          icon: Wrench,
        },
      ],
    },
    {
      label: "OPERATIONS",
      items: [
        {
          label: "Technicians",
          path: "/technicians",
          icon: Users,
        },
        {
          label: "Vehicles",
          path: "/vehicles",
          icon: Truck,
        },
        {
          label: "Parts & Inventory",
          path: "/parts",
          icon: Package,
        },
      ],
    },
    {
      label: "GENERAL",
      items: [notificationItem, settingsItem],
    },
  ],

  TECHNICIAN: [
    {
      label: "MY WORKSPACE",
      items: [
        dashboardItem,
        {
          label: "My Work Orders",
          path: "/my-work-orders",
          icon: Wrench,
        },
        {
          label: "My Schedule",
          path: "/my-schedule",
          icon: CalendarDays,
        },
      ],
    },
    {
      label: "GENERAL",
      items: [notificationItem, settingsItem],
    },
  ],

  CUSTOMER: [
    {
      label: "MY SERVICES",
      items: [
        dashboardItem,
        {
          label: "My Service Requests",
          path: "/service-requests",
          icon: ClipboardList,
        },
        {
          label: "New Service Request",
          path: "/new-service-request",
          icon: PlusCircle,
        },
        {
          label: "Service Tracking",
          path: "/service-tracking",
          icon: Search,
        },
        {
          label: "Service History",
          path: "/service-history",
          icon: History,
        },
      ],
    },
    {
      label: "GENERAL",
      items: [notificationItem, settingsItem],
    },
  ],
};

export default function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const user = useAuthStore((state) => state.user);

  const navigation = user
    ? navigationByRole[user.role]
    : [];

  return (
    <aside className={`sidebar ${collapsed ? "sidebar-collapsed" : ""}`}>
      <div className="sidebar-brand">
        <div className="brand-mark">
          <Wrench size={21} />
        </div>

        {!collapsed && (
          <div className="brand-text">
            <strong>SERVIXA</strong>
            <span>Service Operations</span>
          </div>
        )}

        <button
          type="button"
          className="icon-button sidebar-collapse"
          onClick={() => setCollapsed((current) => !current)}
          aria-label="Toggle sidebar"
          title="Toggle sidebar"
        >
          <ChevronLeft size={17} />
        </button>
      </div>

      <nav className="sidebar-nav">
        {navigation.map((group) => (
          <div className="nav-group" key={group.label}>
            {!collapsed && (
              <div className="nav-label">{group.label}</div>
            )}

            {group.items.map((item) => {
              const Icon = item.icon;

              return (
                <NavLink
                  key={item.path}
                  to={item.path}
                  title={collapsed ? item.label : undefined}
                  className={({ isActive }) =>
                    `nav-item ${isActive ? "nav-item-active" : ""}`
                  }
                >
                  <Icon size={19} strokeWidth={1.8} />
                  {!collapsed && <span>{item.label}</span>}
                </NavLink>
              );
            })}
          </div>
        ))}
      </nav>

      {!collapsed && (
        <div className="sidebar-footer">
          <div className="sidebar-footer-icon">
            <Wrench size={18} />
          </div>

          <div>
            <strong>SERVIXA</strong>
            <span>
              {user ? user.role : "Service Operations"}
            </span>
          </div>
        </div>
      )}
    </aside>
  );
}
