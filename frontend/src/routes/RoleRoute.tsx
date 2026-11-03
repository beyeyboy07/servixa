
import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore, type UserRole } from "../stores/authStore";

interface RoleRouteProps {
  allowedRoles: UserRole[];
}

export default function RoleRoute({
  allowedRoles,
}: RoleRouteProps) {
  const user = useAuthStore((state) => state.user);
  const isAuthenticated = useAuthStore(
    (state) => state.isAuthenticated,
  );

  if (!isAuthenticated || !user) {
    return <Navigate to="/login" replace />;
  }

  if (!allowedRoles.includes(user.role)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
