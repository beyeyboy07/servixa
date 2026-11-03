
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Bell,
  Menu,
  Moon,
  Sun,
  UserRound,
  LogOut,
} from "lucide-react";
import { useAuthStore } from "../../stores/authStore";
import { useNotificationStore } from "../../stores/notificationStore";

interface TopbarProps {
  onMenuClick: () => void;
  isDark: boolean;
  onThemeToggle: () => void;
}

const roleLabels: Record<string, string> = {
  ADMIN: "Administrator",
  DISPATCHER: "Dispatcher",
  TECHNICIAN: "Technician",
  CUSTOMER: "Customer",
};

export default function Topbar({
  onMenuClick,
  isDark,
  onThemeToggle,
}: TopbarProps) {
  const [profileMenuOpen, setProfileMenuOpen] = useState(false);
  const navigate = useNavigate();

  const user = useAuthStore((state) => state.user);
  const isAuthenticated = useAuthStore(
    (state) => state.isAuthenticated,
  );
  const logout = useAuthStore((state) => state.logout);

  const unreadCount = useNotificationStore(
    (state) => state.unreadCount,
  );

  const fetchUnreadCount = useNotificationStore(
    (state) => state.fetchUnreadCount,
  );

  useEffect(() => {
    if (!isAuthenticated) return;

    const refreshUnreadCount = () => {
      if (document.visibilityState === "visible") {
        void fetchUnreadCount();
      }
    };

    refreshUnreadCount();

    const interval = window.setInterval(
      refreshUnreadCount,
      30000,
    );

    document.addEventListener(
      "visibilitychange",
      refreshUnreadCount,
    );

    return () => {
      window.clearInterval(interval);
      document.removeEventListener(
        "visibilitychange",
        refreshUnreadCount,
      );
    };
  }, [isAuthenticated, fetchUnreadCount]);

  function handleLogout() {
    logout();
    useNotificationStore.getState().reset();

    setProfileMenuOpen(false);
    navigate("/login", { replace: true });
  }

  const displayRole = user
    ? roleLabels[user.role] ?? user.role
    : "";

  return (
    <header className="topbar">
      <div className="topbar-left">
        <button
          className="icon-button menu-button"
          onClick={onMenuClick}
          aria-label="Open navigation"
        >
          <Menu size={21} />
        </button>

        <div className="topbar-heading">
          <span className="topbar-eyebrow">WORKSPACE</span>
          <span className="topbar-workspace-name">
            Operations Center
          </span>
        </div>
      </div>

      <div className="topbar-actions">
        <button
          className="icon-button"
          onClick={onThemeToggle}
          aria-label={
            isDark ? "Switch to light theme" : "Switch to dark theme"
          }
          title={isDark ? "Light mode" : "Dark mode"}
        >
          {isDark ? <Sun size={19} /> : <Moon size={19} />}
        </button>

        <button
          className="icon-button notification-button"
          onClick={() => navigate("/notifications")}
          aria-label={`Notifications${
            unreadCount > 0 ? `, ${unreadCount} unread` : ""
          }`}
          title="Notifications"
        >
          <Bell size={19} />

          {unreadCount > 0 && (
            <span className="notification-badge">
              {unreadCount > 99 ? "99+" : unreadCount}
            </span>
          )}
        </button>

        <div className="topbar-divider" />

        <div className="topbar-profile-container">
          <button
            className="user-profile"
            type="button"
            onClick={() =>
              setProfileMenuOpen((current) => !current)
            }
            aria-expanded={profileMenuOpen}
            aria-label="Open user menu"
          >
            <div className="user-avatar">
              <UserRound size={19} />
            </div>

            <div className="user-profile-copy">
              <span className="user-profile-name">
                {user?.name ?? "User"}
              </span>

              <span className="user-profile-role">
                {displayRole}
              </span>
            </div>
          </button>

          {profileMenuOpen && (
            <div className="profile-dropdown">
              <div className="profile-dropdown-info">
                <span className="profile-dropdown-name">
                  {user?.name ?? "User"}
                </span>

                <span className="profile-dropdown-email">
                  {user?.email ?? ""}
                </span>

                <span className="profile-dropdown-role">
                  {displayRole}
                </span>
              </div>

              <button
                className="profile-dropdown-item"
                type="button"
                onClick={handleLogout}
              >
                <LogOut size={17} />
                <span>Logout</span>
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
