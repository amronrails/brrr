import { api, setAccessToken, setRefreshToken, getRefreshToken, clearTokens } from "@/lib/api/client";
import type { AuthResponse, User } from "./types";

function persist(res: AuthResponse) {
  setAccessToken(res.access_token);
  setRefreshToken(res.refresh_token);
}

export async function login(input: { email: string; password: string }): Promise<AuthResponse> {
  const res = await api<AuthResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(input),
    auth: false,
  });
  persist(res);
  return res;
}

export async function register(input: {
  email: string;
  password: string;
  name: string;
}): Promise<AuthResponse> {
  const res = await api<AuthResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify(input),
    auth: false,
  });
  persist(res);
  return res;
}

export async function me(): Promise<User> {
  return api<User>("/auth/me");
}

export async function logout(): Promise<void> {
  const refreshToken = getRefreshToken();
  if (refreshToken) {
    try {
      await api<void>("/auth/logout", {
        method: "POST",
        body: JSON.stringify({ refresh_token: refreshToken }),
        auth: false,
      });
    } catch {
      // ignore network/logout errors; tokens are cleared regardless
    }
  }
  clearTokens();
}
