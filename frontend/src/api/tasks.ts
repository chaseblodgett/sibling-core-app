import { apiClient } from "./client";
import type { Task } from "../types";

export interface CreateTaskPayload {
  title: string;
  description: string;
  assignee_id?: string;
}

export interface UpdateTaskPayload {
  title?: string;
  description?: string;
  status?: Task["status"];
  assignee_id?: string;
}

export const tasksApi = {
  list: (familyId: string) => apiClient.get<Task[]>(`/families/${familyId}/tasks`),
  create: (familyId: string, payload: CreateTaskPayload) =>
    apiClient.post<Task>(`/families/${familyId}/tasks`, payload),
  get: (id: string) => apiClient.get<Task>(`/tasks/${id}`),
  update: (id: string, payload: UpdateTaskPayload) => apiClient.patch<Task>(`/tasks/${id}`, payload),
  assign: (id: string, memberId: string) => apiClient.post<Task>(`/tasks/${id}/assign`, { member_id: memberId }),
  delete: (id: string) => apiClient.delete<void>(`/tasks/${id}`),
};
