import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { HttpError } from "../../../shared/api/http";

export function RegisterPage() {
  const { register } = useAuth();
  const nav = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (username.trim().length < 3) return setError("Username must be at least 3 characters");
    if (password.length < 8) return setError("Password must be at least 8 characters");
    if (password !== confirm) return setError("Passwords do not match");

    setLoading(true);
    try {
      await register(username.trim(), password);
      nav("/", { replace: true });
    } catch (e) {
      if (e instanceof HttpError) {
        setError(`Register failed (${e.status})`);
      } else {
        setError("Register failed");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="container">
      <div className="card">
        <h1>Register</h1>

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
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>

          <label>
            Confirm password
            <input
              type="password"
              autoComplete="new-password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
            />
          </label>

          <button disabled={loading} type="submit">
            {loading ? "Creating…" : "Create account"}
          </button>
        </form>

        <p className="muted">
          Already have an account? <Link to="/login">Login</Link>
        </p>
      </div>
    </div>
  );
}
