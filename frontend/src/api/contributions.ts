import { apiClient } from "./client";
import type { Contribution } from "../types";

export const contributionsApi = {
  list: (familyId: string) => apiClient.get<Contribution[]>(`/families/${familyId}/contributions`),
};
