import { Box, Button, Typography } from "@mui/material";
import type { ReactNode } from "react";

interface EmptyStateProps {
  icon: ReactNode;
  title: string;
  description: string;
  actionLabel?: string;
  onAction?: () => void;
}

// A real empty state (icon + heading + explanation + a way forward)
// instead of a single line of secondary ListItemText — used everywhere
// a list can legitimately be empty, so "nothing here yet" always comes
// with a next step rather than a dead end.
export function EmptyState({ icon, title, description, actionLabel, onAction }: EmptyStateProps) {
  return (
    <Box sx={{ textAlign: "center", py: 6, px: 3 }}>
      <Box sx={{ color: "text.disabled", mb: 1.5, display: "flex", justifyContent: "center", "& svg": { fontSize: 36 } }}>
        {icon}
      </Box>
      <Typography variant="h6" component="h3" gutterBottom>
        {title}
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 360, mx: "auto" }}>
        {description}
      </Typography>
      {actionLabel && onAction && (
        <Button variant="contained" onClick={onAction} sx={{ mt: 2.5 }}>
          {actionLabel}
        </Button>
      )}
    </Box>
  );
}
