import { Box, Paper, Stack, Typography } from "@mui/material";
import type { ReactNode } from "react";
import { FamilyAvatar } from "./FamilyAvatar";
import type { Family } from "../types";
import type { FamilyColor } from "../utils/familyColor";

interface FamilySectionProps {
  family: Family;
  color: FamilyColor;
  /** Right-aligned header content — a stat, a badge, an action button. */
  action?: ReactNode;
  children: ReactNode;
}

// One family's "block": a tinted header band naming the family (DESIGN.md's
// card-feature-tint pattern) plus a matching left border and avatar, so a
// page that stacks several of these reads as clearly separated family
// sections rather than one undifferentiated list.
export function FamilySection({ family, color, action, children }: FamilySectionProps) {
  return (
    <Paper sx={{ borderLeft: `4px solid ${color.border}`, overflow: "hidden" }}>
      <Box
        sx={{
          bgcolor: color.bg,
          px: 3,
          py: 1.5,
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          flexWrap: "wrap",
          gap: 1,
        }}
      >
        <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
          <FamilyAvatar name={family.name} color={color} size={28} />
          <Typography variant="subtitle1" sx={{ fontWeight: 600, color: color.text }}>
            {family.name}
          </Typography>
        </Stack>
        {action}
      </Box>
      {children}
    </Paper>
  );
}
