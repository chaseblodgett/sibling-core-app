import { ApiError } from "./client";
import { tokenStorage } from "./tokenStorage";
import type { Document } from "../types";

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1";

function authHeaders(): HeadersInit {
  const token = tokenStorage.get();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function parseJsonOrThrow<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const body = await response.json().catch(() => ({ message: response.statusText }));
    throw new ApiError(response.status, body.message ?? "Request failed");
  }
  return response.json() as Promise<T>;
}

export const documentsApi = {
  list: (familyId: string, parentId?: string) =>
    fetch(
      `${BASE_URL}/families/${familyId}/documents${parentId ? `?parent_id=${encodeURIComponent(parentId)}` : ""}`,
      { headers: authHeaders() },
    ).then((res) => parseJsonOrThrow<Document[]>(res)),

  get: (id: string) => fetch(`${BASE_URL}/documents/${id}`, { headers: authHeaders() }).then((res) => parseJsonOrThrow<Document>(res)),

  delete: async (id: string) => {
    const res = await fetch(`${BASE_URL}/documents/${id}`, { method: "DELETE", headers: authHeaders() });
    if (!res.ok) {
      const body = await res.json().catch(() => ({ message: res.statusText }));
      throw new ApiError(res.status, body.message ?? "Request failed");
    }
  },

  // Upload is multipart/form-data, so it bypasses the JSON apiClient helper.
  upload: async (familyId: string, file: File, parentId?: string): Promise<Document> => {
    const formData = new FormData();
    formData.append("file", file);
    if (parentId) {
      formData.append("parent_id", parentId);
    }

    const response = await fetch(`${BASE_URL}/families/${familyId}/documents`, {
      method: "POST",
      headers: authHeaders(),
      body: formData,
    });
    return parseJsonOrThrow<Document>(response);
  },

  // Fetches the file with the auth header attached (a plain <a href> can't
  // carry it) and hands the browser a blob to save under its real name.
  download: async (doc: Document) => {
    const response = await fetch(`${BASE_URL}/documents/${doc.id}/download`, { headers: authHeaders() });
    if (!response.ok) {
      throw new ApiError(response.status, "Download failed");
    }
    const blob = await response.blob();
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = doc.file_name;
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
  },
};
