import { createContext, useCallback, useContext, useEffect, useState } from "react";
import type { ReactNode } from "react";
import { ApiError } from "../api/client";
import { familyApi } from "../api/family";
import { useAuth } from "./AuthContext";
import type { Family } from "../types";

interface FamilyContextValue {
  families: Family[];
  loading: boolean;
  error: string | null;
  reload: () => Promise<void>;
}

const FamilyContext = createContext<FamilyContextValue | undefined>(undefined);

export function FamilyProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [families, setFamilies] = useState<Family[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      setFamilies(await familyApi.mine());
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't load your families.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (user) {
      reload();
    } else {
      setFamilies([]);
      setLoading(false);
    }
  }, [user, reload]);

  return <FamilyContext.Provider value={{ families, loading, error, reload }}>{children}</FamilyContext.Provider>;
}

export function useFamilies(): FamilyContextValue {
  const ctx = useContext(FamilyContext);
  if (!ctx) {
    throw new Error("useFamilies must be used within a FamilyProvider");
  }
  return ctx;
}
