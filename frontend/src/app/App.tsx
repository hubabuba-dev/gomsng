import { BrowserRouter, Route, Routes, Link } from "react-router-dom";
import { AuthProvider, useAuth } from "../features/auth/context/AuthContext";
import { ProtectedRoute } from "../features/auth/components/ProtectedRoute";
import { LoginPage } from "../features/auth/pages/LoginPage";
import { RegisterPage } from "../features/auth/pages/RegisterPage";
import { HomePage } from "../features/dashboard/pages/HomePage";

function TopBar() {
  const { authed } = useAuth();

  return (
    <div className="topbar">
      <div className="topbar-inner">
        <Link to="/" className="brand">Messenger</Link>
        <div className="spacer" />
        {!authed ? (
          <div className="row">
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </div>
        ) : (
          <span className="muted">Logged in</span>
        )}
      </div>
    </div>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <TopBar />
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route
            path="/"
            element={
              <ProtectedRoute>
                <HomePage />
              </ProtectedRoute>
            }
          />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
