
import { useState } from "react";
import {
  BrowserRouter,
  Navigate,
  Route,
  Routes,
} from "react-router-dom";

import AppLayout from "./components/AppLayout";
import ProtectedRoute from "./components/ProtectedRoute";
import RoleGuard from "./components/RoleGuard";

import DashboardPage from "./pages/DashboardPage";
import LoginPage from "./pages/LoginPage";
import FeaturePlaceholderPage from "./pages/FeaturePlaceholderPage";
import ServiceRequestsPage from "./pages/ServiceRequestsPage";
import NewServiceRequestPage from "./pages/NewServiceRequestPage";


import type { AuthUser } from "./stores/authStore";

export type Theme = "light" | "dark";

type UserRole = AuthUser["role"];

const ALL_ROLES: UserRole[] = [
  "ADMIN",
  "DISPATCHER",
  "TECHNICIAN",
  "CUSTOMER",
];

const OPERATIONS_ROLES: UserRole[] = ["ADMIN", "DISPATCHER"];
const TECHNICIAN_ROLES: UserRole[] = ["TECHNICIAN"];
const CUSTOMER_ROLES: UserRole[] = ["CUSTOMER"];

export default function App() {
  const [theme, setTheme] = useState<Theme>(() => {
    return localStorage.getItem("servixa-theme") === "dark"
      ? "dark"
      : "light";
  });

  function toggleTheme() {
    setTheme((current) => {
      const next = current === "light" ? "dark" : "light";
      localStorage.setItem("servixa-theme", next);
      return next;
    });
  }

  return (
    <div data-theme={theme}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />

          <Route element={<ProtectedRoute />}>
            <Route
              element={
                <AppLayout
                  theme={theme}
                  onToggleTheme={toggleTheme}
                />
              }
            >
              {/* Semua role */}
              <Route element={<RoleGuard allowedRoles={ALL_ROLES} />}>
                <Route path="/dashboard" element={<DashboardPage />} />

                <Route
                  path="/notifications"
                  element={
                    <FeaturePlaceholderPage title="Notifications" />
                  }
                />

                <Route
                  path="/settings"
                  element={<FeaturePlaceholderPage title="Settings" />}
                />
              </Route>

              {/* Admin dan Dispatcher */}
              <Route
                element={
                  <RoleGuard allowedRoles={OPERATIONS_ROLES} />
                }
              >
                <Route
                  path="/work-orders"
                  element={
                    <FeaturePlaceholderPage title="Work Orders" />
                  }
                />

                <Route
                  path="/work-orders/:id"
                  element={
                    <FeaturePlaceholderPage title="Work Order Details" />
                  }
                />

                <Route
                  path="/technicians"
                  element={
                    <FeaturePlaceholderPage title="Technicians" />
                  }
                />

                <Route
                  path="/vehicles"
                  element={<FeaturePlaceholderPage title="Vehicles" />}
                />

                <Route
                  path="/parts"
                  element={
                    <FeaturePlaceholderPage title="Parts & Inventory" />
                  }
                />
              </Route>

              {/* Technician */}
              <Route
                element={
                  <RoleGuard allowedRoles={TECHNICIAN_ROLES} />
                }
              >
                <Route
                  path="/my-work-orders"
                  element={
                    <FeaturePlaceholderPage title="My Work Orders" />
                  }
                />

                <Route
                  path="/my-schedule"
                  element={
                    <FeaturePlaceholderPage title="My Schedule" />
                  }
                />
              </Route>

              {/* Customer */}
              <Route element={<RoleGuard allowedRoles={CUSTOMER_ROLES} />}>
                <Route
                  path="/service-requests"
                  element={<ServiceRequestsPage />}
                />

                <Route
                  path="/new-service-request"
                  element={<NewServiceRequestPage />}
                />

                <Route
                  path="/my-service-requests"
                  element={
                    <FeaturePlaceholderPage title="My Service Requests" />
                  }
                />

                <Route
                  path="/service-tracking"
                  element={
                    <FeaturePlaceholderPage title="Service Tracking" />
                  }
                />

                <Route
                  path="/service-history"
                  element={
                    <FeaturePlaceholderPage title="Service History" />
                  }
                />
              </Route>
            </Route>
          </Route>

          <Route
            path="/"
            element={<Navigate to="/dashboard" replace />}
          />

          <Route
            path="*"
            element={<Navigate to="/dashboard" replace />}
          />
         
        </Routes>
      </BrowserRouter>
    </div>
  );
}
