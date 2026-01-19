import { Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { ready, authed } = useAuth();

  if (!ready) return <div className="container">Loading…</div>;
  if (!authed) return <Navigate to="/login" replace />;

  return children;
}
