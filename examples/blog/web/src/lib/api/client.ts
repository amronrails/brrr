// A small fetch wrapper that injects the access token, transparently refreshes
// it on a 401 using the stored refresh token, and surfaces API errors in a
// typed form. The access token lives in memory; the refresh token is persisted
// in localStorage.

const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080/api/v1";
const REFRESH_KEY = "blog.refresh_token";

let accessToken: string | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_KEY);
}

export function setRefreshToken(token: string | null) {
  if (token) {
    localStorage.setItem(REFRESH_KEY, token);
  } else {
    localStorage.removeItem(REFRESH_KEY);
  }
}

export function clearTokens() {
  setAccessToken(null);
  setRefreshToken(null);
}

// ApiError carries the HTTP status and any per-field validation messages.
export class ApiError extends Error {
  status: number;
  fields?: Record<string, string>;

  constructor(message: string, status: number, fields?: Record<string, string>) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.fields = fields;
  }
}

interface RequestOptions extends RequestInit {
  auth?: boolean;
  _retry?: boolean;
}

export async function api<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { auth = true, _retry = false, ...init } = options;

  const headers = new Headers(init.headers);
  if (init.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (auth && accessToken) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  }

  const res = await fetch(`${BASE_URL}${path}`, { ...init, headers });

  if (res.status === 401 && auth && !_retry && getRefreshToken()) {
    const refreshed = await tryRefresh();
    if (refreshed) {
      return api<T>(path, { ...options, _retry: true });
    }
  }

  return handle<T>(res);
}

async function tryRefresh(): Promise<boolean> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return false;
  try {
    const res = await fetch(`${BASE_URL}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
    if (!res.ok) {
      clearTokens();
      return false;
    }
    const data = (await res.json()) as { access_token: string; refresh_token: string };
    setAccessToken(data.access_token);
    setRefreshToken(data.refresh_token);
    return true;
  } catch {
    clearTokens();
    return false;
  }
}

async function handle<T>(res: Response): Promise<T> {
  if (res.status === 204) {
    return undefined as T;
  }
  const text = await res.text();
  const data = text ? JSON.parse(text) : null;
  if (!res.ok) {
    const message = (data && data.error) || res.statusText || "request failed";
    throw new ApiError(message, res.status, data?.fields);
  }
  return data as T;
}
