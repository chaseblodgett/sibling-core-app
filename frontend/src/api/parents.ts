import { apiClient } from "./client";
import type { Parent, ParentKeyInfo } from "../types";

export interface CreateParentPayload {
  name: string;
  relation: string;
}

export interface UpdateParentPayload extends Partial<ParentKeyInfo> {
  name?: string;
  relation?: string;
}

export const parentsApi = {
  list: (familyId: string) => apiClient.get<Parent[]>(`/families/${familyId}/parents`),
  create: (familyId: string, payload: CreateParentPayload) =>
    apiClient.post<Parent>(`/families/${familyId}/parents`, payload),
  get: (id: string) => apiClient.get<Parent>(`/parents/${id}`),
  update: (id: string, payload: UpdateParentPayload) => apiClient.patch<Parent>(`/parents/${id}`, payload),
  delete: (id: string) => apiClient.delete<void>(`/parents/${id}`),
};
