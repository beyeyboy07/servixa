
import { Navigate, Route, Routes } from "react-router-dom";

import AppLayout from "../components/layout/AppLayout";
import LoginPage from "../pages/auth/LoginPage";

import DashboardPage from "../pages/dashboard/DashboardPage";
import WorkOrdersPage from "../pages/work-orders/WorkOrdersPage";
import ServiceRequestsPage from "../pages/service-requests/ServiceRequestsPage";
import TechniciansPage from "../pages/technicians/TechniciansPage";
import VehiclesPage from "../pages/vehicles/VehiclesPage";
import PartsPage from "../pages/parts/PartsPage";
import NotificationsPage from "../pages/notifications/NotificationsPage";
import WorkOrderDetailPage from "../pages/work-orders/WorkOrderDetailPage";

import ProtectedRoute from "./ProtectedRoute";
import GuestRoute from "./GuestRoute";
import RoleRoute from "./RoleRoute";

export default function AppRoutes() {
  return (
    <Routes>
      {/* Guest */}
      <Route element={<GuestRoute />}>
        <Route path="/login" element={<LoginPage />} />
      </Route>

      {/* Authenticated Users */}
      <Route element={<ProtectedRoute />}>
        <Route element={<AppLayout />}>
          {/* All Roles */}
          <Route element={<RoleRoute allowedRoles={[
            "ADMIN",
            "DISPATCHER",
            "TECHNICIAN",
            "CUSTOMER",
          ]} />}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/work-orders" element={<WorkOrdersPage />} />
            <Route path="/work-orders/:id" element={<WorkOrderDetailPage />} />
            <Route path="/notifications" element={<NotificationsPage />} />
          </Route>

          {/* Admin, Dispatcher, Customer */}
          <Route element={<RoleRoute allowedRoles={[
            "ADMIN",
            "DISPATCHER",
            "CUSTOMER",
          ]} />}>
            <Route path="/service-requests" element={<ServiceRequestsPage />} />
          </Route>

          {/* Admin, Dispatcher */}
          <Route element={<RoleRoute allowedRoles={[
            "ADMIN",
            "DISPATCHER",
          ]} />}>
            <Route path="/technicians" element={<TechniciansPage />} />
            <Route path="/vehicles" element={<VehiclesPage />} />
            <Route path="/parts" element={<PartsPage />} />
          </Route>
        </Route>
      </Route>

      <Route
        path="*"
        element={<Navigate to="/dashboard" replace />}
      />
    </Routes>
  );
}
