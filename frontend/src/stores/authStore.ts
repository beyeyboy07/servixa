
import { create } from "zustand";

export type UserRole = "ADMIN" | "DISPATCHER" | "TECHNICIAN" | "CUSTOMER";

export interface AuthUser {
  id: number;
  name: string;
  email: string;
  role: UserRole;
}

interface AuthState {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isInitialized: boolean;

  login: (user: AuthUser) => void;
  initialize: () => void;
  logout: () => void;
}

const isValidUser = (value: unknown): value is AuthUser => {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const user = value as Record<string, unknown>;

  return (
    typeof user.id === "number" &&
    typeof user.name === "string" &&
    typeof user.email === "string" &&
    ["ADMIN", "DISPATCHER", "TECHNICIAN", "CUSTOMER"].includes(
      user.role as string,
    )
  );
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isInitialized: false,

  login: (user) => {
    localStorage.setItem("auth_user", JSON.stringify(user));

    set({
      user,
      isAuthenticated: true,
      isInitialized: true,
    });
  },

  initialize: () => {
    try {
      const accessToken = localStorage.getItem("access_token");
      const savedUser = localStorage.getItem("auth_user");

      if (!accessToken || !savedUser) {
        localStorage.removeItem("auth_user");

        set({
          user: null,
          isAuthenticated: false,
          isInitialized: true,
        });

        return;
      }

      const user: unknown = JSON.parse(savedUser);

      if (!isValidUser(user)) {
        throw new Error("Invalid user data");
      }

      set({
        user,
        isAuthenticated: true,
        isInitialized: true,
      });
    } catch (error) {
      console.error("Auth initialization failed:", error);

      try {
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
        localStorage.removeItem("auth_user");
      } finally {
        set({
          user: null,
          isAuthenticated: false,
          isInitialized: true,
        });
      }
    }
  },

  logout: () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("auth_user");

    set({
      user: null,
      isAuthenticated: false,
      isInitialized: true,
    });
  },
}));
