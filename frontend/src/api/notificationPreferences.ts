import { apiClient } from "./client";
import type { NotificationPreferences } from "../types";

export const notificationPreferencesApi = {
  get: () => apiClient.get<NotificationPreferences>("/me/notification-preferences"),
  update: (taskAssignedEmail: boolean) =>
    apiClient.patch<NotificationPreferences>("/me/notification-preferences", {
      task_assigned_email: taskAssignedEmail,
    }),
};
