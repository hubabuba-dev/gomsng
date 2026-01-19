import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import * as authApi from "../api";
import { getAccessToken } from "../../../shared/api/http";

type AuthContextValue = {
  ready: boolean;
  authed: boolean;
  login: (u: string, p: string) => Promise<void>;
  register: (u: string, p: string) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [ready, setReady] = useState(false);
  const [authed, setAuthed] = useState(false);

  // When app starts, try to refresh using cookie (if user already logged in)
  useEffect(() => {
    (async () => {
      try {
        await authApi.refresh();
        setAuthed(Boolean(getAccessToken()));
      } catch {
        setAuthed(false);
      } finally {
        setReady(true);
      }
    })();
  }, []);

  const value = useMemo<AuthContextValue>(() => {
    return {
      ready,
      authed,
      login: async (u, p) => {
        await authApi.login(u, p);
        setAuthed(true);
      },
      register: async (u, p) => {
        await authApi.register(u, p);
        setAuthed(true);
      },
      logout: async () => {
        await authApi.logout();
        setAuthed(false);
      },
    };
  }, [ready, authed]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}
