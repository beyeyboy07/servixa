
import { useState, type FormEvent } from "react";
import { AxiosError } from "axios";
import { Eye, EyeOff, LockKeyhole, Mail, Wrench } from "lucide-react";
import { useLocation, useNavigate } from "react-router-dom";
import { authApi } from "../api/auth";
import { useAuthStore } from "../stores/authStore";

interface ErrorResponse {
  success: false;
  message: string;
  error: string;
}

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();

  const setAuth = useAuthStore((state) => state.setAuth);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setErrorMessage("");

    const normalizedEmail = email.trim();

    if (!normalizedEmail || !password) {
      setErrorMessage("Please enter your email and password.");
      return;
    }

    setLoading(true);

    try {
      const result = await authApi.login({
        email: normalizedEmail,
        password,
      });

      setAuth(
        result.access_token,
        result.refresh_token,
        result.user,
      );

      const destination = (
        location.state as { from?: { pathname?: string } } | null
      )?.from?.pathname;

      navigate(destination || "/dashboard", { replace: true });
    } catch (error) {
      const axiosError = error as AxiosError<ErrorResponse>;

      setErrorMessage(
        axiosError.response?.data?.message ||
          "Unable to sign in. Please check your connection and try again.",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="login-page">
      <section className="login-brand-panel">
        <div className="login-brand">
          <div className="login-brand-mark">
            <Wrench size={24} />
          </div>
          <span>SERVIXA</span>
        </div>

        <div className="login-brand-content">
          <span className="login-overline">SERVICE OPERATIONS PLATFORM</span>
          <h1>
            Keep your
            <br />
            operations
            <br />
            <span>moving forward.</span>
          </h1>
          <p>
            Manage service requests, coordinate technicians, and keep
            every work order moving in one place.
          </p>
        </div>

        <div className="login-brand-footer">
          Fleet &amp; Field Service Management
        </div>
      </section>

      <section className="login-form-panel">
        <div className="login-form-container">
          <div className="login-mobile-brand">
            <div className="login-brand-mark">
              <Wrench size={22} />
            </div>
            <span>SERVIXA</span>
          </div>

          <div className="login-heading">
            <span className="login-overline">WELCOME BACK</span>
            <h2>Sign in to your account</h2>
            <p>Enter your credentials to access your workspace.</p>
          </div>

          <form className="login-form" onSubmit={handleSubmit}>
            {errorMessage && (
              <div className="login-error" role="alert">
                {errorMessage}
              </div>
            )}

            <div className="login-field">
              <label htmlFor="email">Email address</label>
              <div className="login-input-wrapper">
                <Mail size={18} />
                <input
                  id="email"
                  type="email"
                  autoComplete="email"
                  placeholder="you@company.com"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  required
                />
              </div>
            </div>

            <div className="login-field">
              <div className="login-password-label">
                <label htmlFor="password">Password</label>
              </div>

              <div className="login-input-wrapper">
                <LockKeyhole size={18} />
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="current-password"
                  placeholder="Enter your password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  required
                />

                <button
                  className="password-toggle"
                  type="button"
                  onClick={() => setShowPassword((current) => !current)}
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>

            <button
              className="login-submit"
              type="submit"
              disabled={loading}
            >
              {loading ? "Signing in..." : "Sign in"}
            </button>
          </form>

          <p className="login-security-note">
            Secure access to your service operations workspace.
          </p>
        </div>

        <div className="login-copyright">
          © {new Date().getFullYear()} SERVIXA. All rights reserved.
        </div>
      </section>
    </main>
  );
}
