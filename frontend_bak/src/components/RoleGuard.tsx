
import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore, type AuthUser } from "../stores/authStore";

type UserRole = AuthUser["role"];

interface RoleGuardProps {
  allowedRoles: UserRole[];
}

export default function RoleGuard({
  allowedRoles,
}: RoleGuardProps) {
  const user = useAuthStore((state) => state.user);

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (!allowedRoles.includes(user.role)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
