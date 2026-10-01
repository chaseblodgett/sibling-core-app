import { apiClient } from "./client";
import type { Balance, Expense } from "../types";

export interface ExpenseSplitPayload {
  member_id: string;
  share_cents: number;
}

export interface CreateExpensePayload {
  description: string;
  amount_cents: number;
  currency: string;
  category: string;
  splits: ExpenseSplitPayload[];
}

export const expensesApi = {
  list: (familyId: string) => apiClient.get<Expense[]>(`/families/${familyId}/expenses`),
  create: (familyId: string, payload: CreateExpensePayload) =>
    apiClient.post<Expense>(`/families/${familyId}/expenses`, payload),
  get: (id: string) => apiClient.get<Expense>(`/expenses/${id}`),
  delete: (id: string) => apiClient.delete<void>(`/expenses/${id}`),
  settleSplit: (splitId: string) => apiClient.post<void>(`/expenses/splits/${splitId}/settle`),
  balances: (familyId: string) => apiClient.get<Balance[]>(`/families/${familyId}/balances`),
};
