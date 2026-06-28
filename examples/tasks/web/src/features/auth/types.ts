export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  created_at: string;
}

export interface AuthResponse {
  user: User;
  access_token: string;
  expires_at: string;
  refresh_token: string;
}
