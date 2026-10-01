import { apiClient } from "./client";
import type { Document } from "../types";

export const documentsApi = {
  list: (familyId: string) => apiClient.get<Document[]>(`/families/${familyId}/documents`),
  get: (id: string) => apiClient.get<Document>(`/documents/${id}`),
  delete: (id: string) => apiClient.delete<void>(`/documents/${id}`),
  // Upload is multipart/form-data, so it bypasses the JSON apiClient helper.
  upload: async (familyId: string, file: File): Promise<Document> => {
    const formData = new FormData();
    formData.append("file", file);

    const baseUrl = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1";
    const response = await fetch(`${baseUrl}/families/${familyId}/documents`, {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.statusText}`);
    }

    return response.json();
  },
};
