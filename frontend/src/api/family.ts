import { apiClient } from "./client";
import type { Family, FamilyInvite, Member } from "../types";

export const familyApi = {
  create: (name: string) => apiClient.post<Family>("/families", { name }),
  get: (id: string) => apiClient.get<Family>(`/families/${id}`),
  join: (inviteCode: string) => apiClient.post<Family>("/families/join", { invite_code: inviteCode }),
  invite: (familyId: string, email: string) =>
    apiClient.post<FamilyInvite>(`/families/${familyId}/invites`, { email }),
  listMembers: (familyId: string) => apiClient.get<Member[]>(`/families/${familyId}/members`),
  removeMember: (familyId: string, memberId: string) =>
    apiClient.delete<void>(`/families/${familyId}/members/${memberId}`),
};
