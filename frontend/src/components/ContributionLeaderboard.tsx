import AssignmentTurnedInIcon from "@mui/icons-material/AssignmentTurnedIn";
import AttachMoneyIcon from "@mui/icons-material/AttachMoney";
import { Stack, Typography } from "@mui/material";
import type { Contribution } from "../types";

function formatCents(cents: number) {
  return (cents / 100).toLocaleString(undefined, { style: "currency", currency: "USD" });
}

interface ContributionLeaderboardProps {
  entries: Contribution[];
}

// Who's picking up the most tasks, who's fronting the most expenses —
// a transparency tool, not a scoreboard, so everyone shows up even at
// zero. Ranked by tasks completed (the plan's primary signal), then by
// amount fronted — hiding low contributors would defeat the point of a
// feature meant to surface imbalance.
export function ContributionLeaderboard({ entries }: ContributionLeaderboardProps) {
  if (entries.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ px: 2, py: 1.5 }}>
        Nothing logged yet.
      </Typography>
    );
  }

  return (
    <Stack spacing={1} sx={{ py: 0.5 }}>
      {entries.map((entry) => (
        <Stack
          key={entry.member_id}
          direction="row"
          spacing={2}
          sx={{ alignItems: "center", justifyContent: "space-between", px: 2, py: 0.5 }}
        >
          <Typography variant="body2" sx={{ fontWeight: 500 }}>
            {entry.member_name}
          </Typography>
          <Stack direction="row" spacing={2.5}>
            <Stack direction="row" spacing={0.5} sx={{ alignItems: "center" }}>
              <AssignmentTurnedInIcon fontSize="inherit" sx={{ color: "text.disabled" }} />
              <Typography variant="body2" color="text.secondary">
                {entry.tasks_completed} {entry.tasks_completed === 1 ? "task" : "tasks"}
              </Typography>
            </Stack>
            <Stack direction="row" spacing={0.5} sx={{ alignItems: "center" }}>
              <AttachMoneyIcon fontSize="inherit" sx={{ color: "text.disabled" }} />
              <Typography variant="body2" color="text.secondary">
                {entry.expenses_paid_count > 0 ? formatCents(entry.expenses_paid_cents) : "$0"}
              </Typography>
            </Stack>
          </Stack>
        </Stack>
      ))}
    </Stack>
  );
}
