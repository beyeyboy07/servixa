
import { Outlet } from "react-router-dom";
import type { Theme } from "../App";
import Sidebar from "./Sidebar";
import Header from "./Header";

interface AppLayoutProps {
  theme: Theme;
  onToggleTheme: () => void;
}

export default function AppLayout({
  theme,
  onToggleTheme,
}: AppLayoutProps) {
  return (
    <div className="app-shell">
      <Sidebar />

      <div className="main-area">
        <Header theme={theme} onToggleTheme={onToggleTheme} />

        <main className="page-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
