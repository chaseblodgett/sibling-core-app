import { Alert, Snackbar } from "@mui/material";
import { createContext, useCallback, useContext, useState } from "react";
import type { ReactNode } from "react";

type Severity = "success" | "error" | "info";

interface NotificationContextValue {
  /** Shows a brief toast. Defaults to "success" styling. */
  notify: (message: string, severity?: Severity) => void;
}

const NotificationContext = createContext<NotificationContextValue | undefined>(undefined);

// A single app-wide toast queue — simpler than per-page Snackbar state,
// and means every page gets consistent "did that work?" feedback for
// create/update/delete actions for free.
export function NotificationProvider({ children }: { children: ReactNode }) {
  const [open, setOpen] = useState(false);
  const [message, setMessage] = useState("");
  const [severity, setSeverity] = useState<Severity>("success");

  const notify = useCallback((msg: string, sev: Severity = "success") => {
    setMessage(msg);
    setSeverity(sev);
    setOpen(true);
  }, []);

  return (
    <NotificationContext.Provider value={{ notify }}>
      {children}
      <Snackbar
        open={open}
        autoHideDuration={3500}
        onClose={() => setOpen(false)}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      >
        <Alert onClose={() => setOpen(false)} severity={severity} variant="filled" sx={{ width: "100%" }}>
          {message}
        </Alert>
      </Snackbar>
    </NotificationContext.Provider>
  );
}

export function useNotify(): NotificationContextValue["notify"] {
  const ctx = useContext(NotificationContext);
  if (!ctx) {
    throw new Error("useNotify must be used within a NotificationProvider");
  }
  return ctx.notify;
}
