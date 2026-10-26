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
  google: (idToken: string) => apiClient.post<AuthResponse>("/auth/google", { id_token: idToken }),
  me: () => apiClient.get<AuthUser>("/auth/me"),
};
