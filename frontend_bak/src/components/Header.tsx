
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Bell,
  LogOut,
  Moon,
  Search,
  Sun,
} from "lucide-react";

import type { Theme } from "../App";
import { authApi } from "../api/auth";
import { useAuthStore } from "../stores/authStore";

interface HeaderProps {
  theme: Theme;
  onToggleTheme: () => void;
}

export default function Header({
  theme,
  onToggleTheme,
}: HeaderProps) {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const refreshToken = useAuthStore((state) => state.refreshToken);
  const clearAuth = useAuthStore((state) => state.clearAuth);

  const [loggingOut, setLoggingOut] = useState(false);

  async function handleLogout() {
    if (loggingOut) return;

    setLoggingOut(true);

    try {
      if (refreshToken) {
        await authApi.logout(refreshToken);
      }
    } catch {
      // Clear local session even if the logout request fails.
    } finally {
      clearAuth();
      navigate("/login", { replace: true });
      setLoggingOut(false);
    }
  }

  return (
    <header className="topbar">
      <div className="topbar-search">
        <Search size={18} />
        <input
          type="search"
          placeholder="Search anything..."
          aria-label="Search"
        />
      </div>

      <div className="topbar-actions">
        <button
          className="icon-button theme-toggle"
          onClick={onToggleTheme}
          aria-label={`Switch to ${theme === "light" ? "dark" : "light"} mode`}
          title="Toggle theme"
        >
          {theme === "light" ? <Moon size={19} /> : <Sun size={19} />}
        </button>

        <button
          className="icon-button notification-button"
          aria-label="Notifications"
          title="Notifications"
          onClick={() => navigate("/notifications")}
        >
          <Bell size={19} />
          <span className="notification-dot" />
        </button>

        <div className="topbar-divider" />

        <div className="user-profile">
          <div className="user-avatar">
            {user?.name
              ?.split(/\s+/)
              .slice(0, 2)
              .map((part) => part[0])
              .join("")
              .toUpperCase() || "U"}
          </div>

          <div className="user-details">
            <strong>{user?.name || "User"}</strong>
            <span>{user?.role || "Account"}</span>
          </div>
        </div>

        <button
          className="icon-button logout-button"
          onClick={handleLogout}
          disabled={loggingOut}
          aria-label="Sign out"
          title={loggingOut ? "Signing out..." : "Sign out"}
        >
          <LogOut size={18} />
        </button>
      </div>
    </header>
  );
}
