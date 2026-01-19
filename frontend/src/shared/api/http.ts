const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string;

type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export class HttpError extends Error {
  public status: number;
  public details?: unknown;

  constructor(message: string, status: number, details?: unknown) {
    super(message);
    this.name = "HttpError";
    this.status = status;
    this.details = details;
  }
}


let accessToken: string | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getAccessToken() {
  return accessToken;
}

async function parseJsonSafe(res: Response) {
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

export async function http<T>(
  path: string,
  opts: {
    method?: HttpMethod;
    body?: unknown;
    auth?: boolean; // attach Authorization header
    retryOn401?: boolean; // try refresh once
  } = {}
): Promise<T> {
  const { method = "GET", body, auth = false, retryOn401 = true } = opts;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (auth && accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  const res = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
    credentials: "include", // IMPORTANT: sends refresh cookie
  });

  if (res.ok) {
    return (await parseJsonSafe(res)) as T;
  }

  // If access token expired, try refresh once then retry original request
  if (res.status === 401 && retryOn401) {
    try {
      const refreshed = await fetch(`${API_BASE_URL}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });

      if (!refreshed.ok) throw new Error("Refresh failed");

      const data = (await parseJsonSafe(refreshed)) as { accessToken: string };
      setAccessToken(data.accessToken);

      return http<T>(path, { method, body, auth, retryOn401: false });
    } catch {
      // fall through and throw error below
    }
  }

  const details = await parseJsonSafe(res);
  throw new HttpError("Request failed", res.status, details);
}