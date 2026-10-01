import { apiClient } from "./client";
import type { AuthUser } from "../types";

export interface SignUpPayload {
  email: string;
  password: string;
  full_name: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export const authApi = {
  signUp: (payload: SignUpPayload) => apiClient.post<AuthUser>("/auth/signup", payload),
  login: (payload: LoginPayload) => apiClient.post<AuthUser>("/auth/login", payload),
  logout: () => apiClient.post<void>("/auth/logout"),
  me: () => apiClient.get<AuthUser>("/auth/me"),
};
