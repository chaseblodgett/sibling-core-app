import { ApiError } from "./client";
import { tokenStorage } from "./tokenStorage";

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1";

function authHeaders(): HeadersInit {
  const token = tokenStorage.get();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

// Fetches a CSV export with the auth header attached (a plain <a href>
// can't carry it) and hands the browser a blob to save — same pattern as
// documentsApi.download.
async function downloadCSV(url: string, filename: string) {
  const response = await fetch(url, { headers: authHeaders() });
  if (!response.ok) {
    const body = await response.json().catch(() => ({ message: response.statusText }));
    throw new ApiError(response.status, body.message ?? "Export failed");
  }
  const blob = await response.blob();
  const objectUrl = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = objectUrl;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(objectUrl);
}

export const exportApi = {
  expenses: (familyId: string) => downloadCSV(`${BASE_URL}/families/${familyId}/export/expenses.csv`, "expenses.csv"),
  tasks: (familyId: string) => downloadCSV(`${BASE_URL}/families/${familyId}/export/tasks.csv`, "tasks.csv"),
};
