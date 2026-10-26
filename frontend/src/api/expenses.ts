import { apiClient } from "./client";
import type { Balance, Expense, ExpenseSplit } from "../types";

export interface ExpenseSplitPayload {
  member_id: string;
  share_cents: number;
}

export interface CreateExpensePayload {
  description: string;
  amount_cents: number;
  currency: string;
  category: string;
  parent_id?: string;
  paid_by?: string;
  splits: ExpenseSplitPayload[];
}

export const expensesApi = {
  list: (familyId: string, parentId?: string) =>
    apiClient.get<Expense[]>(
      `/families/${familyId}/expenses${parentId ? `?parent_id=${encodeURIComponent(parentId)}` : ""}`,
    ),
  create: (familyId: string, payload: CreateExpensePayload) =>
    apiClient.post<Expense>(`/families/${familyId}/expenses`, payload),
  get: (id: string) => apiClient.get<Expense>(`/expenses/${id}`),
  delete: (id: string) => apiClient.delete<void>(`/expenses/${id}`),
  listSplits: (expenseId: string) => apiClient.get<ExpenseSplit[]>(`/expenses/${expenseId}/splits`),
  settleSplit: (splitId: string) => apiClient.post<ExpenseSplit>(`/expenses/splits/${splitId}/settle`),
  balances: (familyId: string) => apiClient.get<Balance[]>(`/families/${familyId}/balances`),
};
