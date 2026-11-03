
import { useState } from "react";
import type { FormEvent } from "react";
import {
  Eye,
  EyeOff,
  LockKeyhole,
  Mail,
  Wrench,
} from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";
import { useAuthStore, type AuthUser } from "../../stores/authStore";
import api from "../../services/api";

interface LoginResponse {
  success: boolean;
  message: string;
  data: {
    access_token: string;
    refresh_token: string;
    user: {
      id: number;
      name: string;
      email: string;
      role: AuthUser["role"];
      is_active: boolean;
    };
  };
}

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const navigate = useNavigate();
  const location = useLocation();
  const login = useAuthStore((state) => state.login);

  const redirectTo =
    (location.state as { from?: { pathname?: string } } | null)
      ?.from?.pathname ?? "/dashboard";

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (loading) return;

    setLoading(true);
    setError("");

    try {
      const response = await api.post<LoginResponse>("/auth/login", {
        email,
        password,
      });

      const result = response.data;

      if (!result.success || !result.data?.access_token || !result.data.user) {
        throw new Error(result.message || "Login failed.");
      }

      if (!result.data.user.is_active) {
        setError("Your account is inactive. Please contact the administrator.");
        return;
      }

      localStorage.setItem("access_token", result.data.access_token);
      localStorage.setItem("refresh_token", result.data.refresh_token);

      const user: AuthUser = {
        id: result.data.user.id,
        name: result.data.user.name,
        email: result.data.user.email,
        role: result.data.user.role,
      };

      login(user);
      navigate(redirectTo, { replace: true });
    } catch (err: unknown) {
      const message =
        typeof err === "object" &&
        err !== null &&
        "response" in err
          ? (
              err as {
                response?: {
                  data?: { message?: string };
                };
              }
            ).response?.data?.message
          : undefined;

      setError(
        message ||
          (err instanceof Error && !("response" in err)
            ? err.message
            : "Unable to sign in. Check your credentials and try again."),
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="login-page">
      <section className="login-brand-panel">
        <div className="login-brand">
          <div className="login-brand-icon">
            <Wrench size={24} />
          </div>
          <span>SERVIXA</span>
        </div>

        <div className="login-brand-content">
          <span className="login-eyebrow">FLEET & FIELD SERVICE</span>
          <h1>Keep your operations moving.</h1>
          <p>
            Manage service requests, coordinate technicians, and keep your
            fleet running from one workspace.
          </p>
        </div>

        <div className="login-brand-footer">
          Fleet management made simpler.
        </div>
      </section>

      <section className="login-form-panel">
        <div className="login-form-container">
          <div className="login-mobile-brand">
            <div className="login-brand-icon">
              <Wrench size={22} />
            </div>
            <span>SERVIXA</span>
          </div>

          <div className="login-heading">
            <h2>Welcome back</h2>
            <p>Enter your account details to sign in.</p>
          </div>

          <form className="login-form" onSubmit={handleSubmit}>
            <label htmlFor="login-email">Email address</label>
            <div className="login-input-wrapper">
              <Mail size={18} />
              <input
                id="login-email"
                type="email"
                placeholder="you@company.com"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                autoComplete="username"
                required
                disabled={loading}
              />
            </div>

            <label htmlFor="login-password">Password</label>
            <div className="login-input-wrapper">
              <LockKeyhole size={18} />
              <input
                id="login-password"
                type={showPassword ? "text" : "password"}
                placeholder="Enter your password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                autoComplete="current-password"
                required
                disabled={loading}
              />
              <button
                className="login-password-toggle"
                type="button"
                onClick={() => setShowPassword((current) => !current)}
                aria-label={showPassword ? "Hide password" : "Show password"}
                disabled={loading}
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>

            {error && (
              <p className="login-error" role="alert">
                {error}
              </p>
            )}

            <button
              className="login-submit-button"
              type="submit"
              disabled={loading}
            >
              {loading ? "Signing in..." : "Sign In"}
            </button>
          </form>

          <p className="login-copyright">
            © 2026 SERVIXA. Fleet & Field Service Management.
          </p>
        </div>
      </section>
    </main>
  );
}
