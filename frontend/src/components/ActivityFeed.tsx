import { List, ListItem, ListItemText, Typography } from "@mui/material";
import type { ActivityEntry, ActivityVerb } from "../types";

const VERB_PHRASE: Record<ActivityVerb, (label: string) => string> = {
  "task.created": (label) => `added a task: ${label}`,
  "task.completed": (label) => `completed ${label}`,
  "expense.created": (label) => `logged an expense: ${label}`,
  "document.uploaded": (label) => `uploaded ${label}`,
  "parent.created": (label) => `added ${label} to the family`,
  "member.joined": () => "joined the family",
};

function describe(entry: ActivityEntry): string {
  const phrase = VERB_PHRASE[entry.verb];
  return phrase ? phrase(entry.subject_label) : "did something";
}

function timeAgo(iso: string): string {
  const seconds = Math.max(0, (Date.now() - new Date(iso).getTime()) / 1000);
  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d ago`;
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric" });
}

interface ActivityFeedProps {
  entries: ActivityEntry[];
}

// A reverse-chronological "what's been happening" list — answers "what
// did I miss" for siblings coordinating async, rather than only ever
// showing current-state stats. Entries are self-contained (actor_name,
// subject_label already resolved server-side), so this never needs to
// cross-reference members/tasks/parents data to render.
export function ActivityFeed({ entries }: ActivityFeedProps) {
  if (entries.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ px: 2, py: 1.5 }}>
        No activity yet.
      </Typography>
    );
  }

  return (
    <List dense sx={{ py: 0.5 }}>
      {entries.map((entry) => (
        <ListItem key={entry.id} sx={{ py: 0.5 }}>
          <ListItemText
            primary={
              <>
                <Typography component="span" variant="body2" sx={{ fontWeight: 500 }}>
                  {entry.actor_name}
                </Typography>{" "}
                <Typography component="span" variant="body2" color="text.secondary">
                  {describe(entry)}
                </Typography>
              </>
            }
            secondary={timeAgo(entry.created_at)}
          />
        </ListItem>
      ))}
    </List>
  );
}
