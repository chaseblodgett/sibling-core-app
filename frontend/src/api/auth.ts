import { apiClient } from "./client";
import type { AuthResponse, AuthUser } from "../types";

export interface SignUpPayload {
  email: string;
  password: string;
  name: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export const authApi = {
  signUp: (payload: SignUpPayload) => apiClient.post<AuthResponse>("/auth/signup", payload),
  login: (payload: LoginPayload) => apiClient.post<AuthResponse>("/auth/login", payload),
  me: () => apiClient.get<AuthUser>("/auth/me"),
};
