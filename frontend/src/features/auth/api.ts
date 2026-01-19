import { http, setAccessToken } from "../../shared/api/http";

export type AuthResponse = { accessToken: string };

export async function login(username: string, password: string) {
  const data = await http<AuthResponse>("/auth/login", {
    method: "POST",
    body: { username, password },
  });
  setAccessToken(data.accessToken);
}

export async function register(username: string, password: string) {
  const data = await http<AuthResponse>("/auth/register", {
    method: "POST",
    body: { username, password },
  });
  setAccessToken(data.accessToken);
}

export async function refresh() {
  const data = await http<AuthResponse>("/auth/refresh", {
    method: "POST",
    retryOn401: false,
  });
  setAccessToken(data.accessToken);
}

export async function logout() {
  await http<void>("/auth/logout", { method: "POST" });
  setAccessToken(null);
}
