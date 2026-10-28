import { apiClient } from "./client";
import type { RecurrenceRule, Task } from "../types";

export interface CreateTaskPayload {
  title: string;
  description: string;
  assignee_id?: string;
  parent_id?: string;
  due_date?: string;
  recurrence_rule?: RecurrenceRule;
}

export interface UpdateTaskPayload {
  title?: string;
  description?: string;
  status?: Task["status"];
  assignee_id?: string;
  parent_id?: string;
}

export const tasksApi = {
  list: (familyId: string, parentId?: string) =>
    apiClient.get<Task[]>(
      `/families/${familyId}/tasks${parentId ? `?parent_id=${encodeURIComponent(parentId)}` : ""}`,
    ),
  create: (familyId: string, payload: CreateTaskPayload) =>
    apiClient.post<Task>(`/families/${familyId}/tasks`, payload),
  get: (id: string) => apiClient.get<Task>(`/tasks/${id}`),
  update: (id: string, payload: UpdateTaskPayload) => apiClient.patch<Task>(`/tasks/${id}`, payload),
  assign: (id: string, memberId: string) => apiClient.post<Task>(`/tasks/${id}/assign`, { member_id: memberId }),
  delete: (id: string) => apiClient.delete<void>(`/tasks/${id}`),
};
