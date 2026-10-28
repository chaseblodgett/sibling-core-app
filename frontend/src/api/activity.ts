import { apiClient } from "./client";
import type { ActivityEntry } from "../types";

export const activityApi = {
  list: (familyId: string, limit = 10) =>
    apiClient.get<ActivityEntry[]>(`/families/${familyId}/activity?limit=${limit}`),
};
