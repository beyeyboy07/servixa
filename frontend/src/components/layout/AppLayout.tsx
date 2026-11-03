
import { useEffect, useState } from "react";
import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import Topbar from "./Topbar";

const THEME_KEY = "servixa_theme";

export default function AppLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const [isDark, setIsDark] = useState(() => {
    const savedTheme = localStorage.getItem(THEME_KEY);

    if (savedTheme === "dark") return true;
    if (savedTheme === "light") return false;

    return document.documentElement.dataset.theme === "dark";
  });

  useEffect(() => {
    const theme = isDark ? "dark" : "light";

    document.documentElement.dataset.theme = theme;
    localStorage.setItem(THEME_KEY, theme);
  }, [isDark]);

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth > 1024) {
        setSidebarOpen(false);
      }
    };

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return (
    <div className="app-shell">
      <Sidebar
        open={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
      />

      <div className="app-main">
        <Topbar
          onMenuClick={() => setSidebarOpen(true)}
          isDark={isDark}
          onThemeToggle={() => setIsDark((current) => !current)}
        />

        <main className="main-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
