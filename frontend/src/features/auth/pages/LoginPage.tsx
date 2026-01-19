import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { HttpError } from "../../../shared/api/http";

export function LoginPage() {
  const { login } = useAuth();
  const nav = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    // tiny client validation (keep it simple)
    if (username.trim().length < 3) return setError("Username must be at least 3 characters");
    if (password.length < 8) return setError("Password must be at least 8 characters");

    setLoading(true);
    try {
      await login(username.trim(), password);
      nav("/", { replace: true });
    } catch (e) {
      if (e instanceof HttpError) {
        setError(`Login failed (${e.status})`);
      } else {
        setError("Login failed");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="container">
      <div className="card">
        <h1>Login</h1>

        {error && <div className="alert">{error}</div>}

        <form className="form" onSubmit={onSubmit}>
          <label>
            Username
            <input
              autoComplete="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </label>

          <label>
            Password
            <input
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>

          <button disabled={loading} type="submit">
            {loading ? "Signing in…" : "Sign in"}
          </button>
        </form>

        <p className="muted">
          No account? <Link to="/register">Register</Link>
        </p>
      </div>
    </div>
  );
}
