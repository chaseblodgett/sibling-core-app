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
  // null means "All families" — the aggregated view. Only meaningful (and
  // only ever shown as a choice in the UI) when there are 2+ families;
  // with exactly one, callers can treat this as that family's id anyway.
  selectedFamilyId: string | null;
  setSelectedFamilyId: (id: string | null) => void;
  reload: () => Promise<void>;
}

const FamilyContext = createContext<FamilyContextValue | undefined>(undefined);

export function FamilyProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [families, setFamilies] = useState<Family[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedFamilyId, setSelectedFamilyId] = useState<string | null>(null);

  const reload = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const mine = await familyApi.mine();
      setFamilies(mine);
      // Default to "All families" once there's more than one; with zero or
      // one, there's nothing to choose between, so pin it directly.
      setSelectedFamilyId(mine.length === 1 ? mine[0].id : null);
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
      setSelectedFamilyId(null);
      setLoading(false);
    }
  }, [user, reload]);

  return (
    <FamilyContext.Provider value={{ families, loading, error, selectedFamilyId, setSelectedFamilyId, reload }}>
      {children}
    </FamilyContext.Provider>
  );
}

export function useFamilies(): FamilyContextValue {
  const ctx = useContext(FamilyContext);
  if (!ctx) {
    throw new Error("useFamilies must be used within a FamilyProvider");
  }
  return ctx;
}
