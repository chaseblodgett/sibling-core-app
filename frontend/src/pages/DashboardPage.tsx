import AddIcon from "@mui/icons-material/Add";
import AssignmentIcon from "@mui/icons-material/Assignment";
import AttachMoneyIcon from "@mui/icons-material/AttachMoney";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import DescriptionIcon from "@mui/icons-material/Description";
import GroupAddIcon from "@mui/icons-material/GroupAdd";
import PersonAddAlt1Icon from "@mui/icons-material/PersonAddAlt1";
import RadioButtonUncheckedIcon from "@mui/icons-material/RadioButtonUnchecked";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Divider,
  Grid,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Paper,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { type FormEvent, type ReactNode, useEffect, useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import { activityApi } from "../api/activity";
import { ApiError } from "../api/client";
import { contributionsApi } from "../api/contributions";
import { documentsApi } from "../api/documents";
import { expensesApi } from "../api/expenses";
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { tasksApi } from "../api/tasks";
import { ActivityFeed } from "../components/ActivityFeed";
import { ContributionLeaderboard } from "../components/ContributionLeaderboard";
import { FamilySection } from "../components/FamilySection";
import { useAuth } from "../context/AuthContext";
import { useFamilies } from "../context/FamilyContext";
import { useNotify } from "../context/NotificationContext";
import { familyColorMap } from "../utils/familyColor";
import type { FamilyColor } from "../utils/familyColor";
import type { ActivityEntry, Contribution, Family } from "../types";

interface FamilyStats {
  openTasks: number;
  balanceCents: number | null; // null = nothing owed either way
  documents: number;
  parents: number;
}

function formatCents(cents: number) {
  return (cents / 100).toLocaleString(undefined, { style: "currency", currency: "USD" });
}

export function DashboardPage() {
  const { user } = useAuth();
  const { families, loading: familiesLoading, reload } = useFamilies();
  const notify = useNotify();

  const [stats, setStats] = useState<Record<string, FamilyStats>>({});
  const [activity, setActivity] = useState<Record<string, ActivityEntry[]>>({});
  const [contributions, setContributions] = useState<Record<string, Contribution[]>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [createName, setCreateName] = useState("");
  const [joinCode, setJoinCode] = useState("");
  const [onboardSubmitting, setOnboardSubmitting] = useState(false);

  useEffect(() => {
    if (familiesLoading) return;
    if (families.length === 0) {
      setLoading(false);
      return;
    }

    let cancelled = false;
    setLoading(true);
    setError(null);

    Promise.all(
      families.map(async (family) => {
        const [tasks, documents, members, balances, parents, recentActivity, memberContributions] = await Promise.all([
          tasksApi.list(family.id),
          documentsApi.list(family.id),
          familyApi.listMembers(family.id),
          expensesApi.balances(family.id),
          parentsApi.list(family.id),
          activityApi.list(family.id, 5),
          contributionsApi.list(family.id),
        ]);
        const myMemberId = members.find((m) => m.user_id === user?.id)?.id;
        const myBalance = balances.find((b) => b.member_id === myMemberId);
        return {
          familyId: family.id,
          stats: {
            openTasks: tasks.filter((t) => t.status !== "done").length,
            balanceCents: myBalance?.net_cents ?? null,
            documents: documents.length,
            parents: parents.length,
          } satisfies FamilyStats,
          activity: recentActivity,
          contributions: memberContributions,
        };
      }),
    )
      .then((entries) => {
        if (!cancelled) {
          setStats(Object.fromEntries(entries.map((e) => [e.familyId, e.stats])));
          setActivity(Object.fromEntries(entries.map((e) => [e.familyId, e.activity])));
          setContributions(Object.fromEntries(entries.map((e) => [e.familyId, e.contributions])));
        }
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof ApiError ? err.message : "Couldn't load your dashboard.");
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [families, familiesLoading, user?.id]);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setOnboardSubmitting(true);
    try {
      await familyApi.create(createName);
      setCreateName("");
      notify(`${createName} is ready to go.`);
      await reload();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't create the family.");
    } finally {
      setOnboardSubmitting(false);
    }
  }

  async function handleJoin(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setOnboardSubmitting(true);
    try {
      await familyApi.join(joinCode);
      setJoinCode("");
      notify("You're in!");
      await reload();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't join with that invite code.");
    } finally {
      setOnboardSubmitting(false);
    }
  }

  if (familiesLoading || loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 8 }}>
        <CircularProgress sx={{ color: "primary.main" }} />
      </Box>
    );
  }

  // Brand-new user, zero families: this IS the start of the app, so the
  // start action lives right here rather than behind a link to /family.
  if (families.length === 0) {
    return (
      <>
        <Typography variant="h4" component="h1" gutterBottom>
          Welcome{user ? `, ${user.name}` : ""}
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ mb: 4, maxWidth: 560 }}>
          Caretree helps siblings coordinate caring for a parent together — shared tasks, split expenses, and
          documents everyone can find. Start by creating a family, or join one with a code from a sibling.
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Grid container spacing={2}>
          <Grid size={{ xs: 12, sm: 6 }}>
            <Paper sx={{ p: 3, height: "100%" }}>
              <Stack spacing={2} component="form" onSubmit={handleCreate} sx={{ height: "100%" }}>
                <Box sx={{ flexGrow: 1 }}>
                  <Typography variant="h6" component="h2">
                    Create a family
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    You'll be the owner — invite siblings once it's set up.
                  </Typography>
                </Box>
                <TextField
                  label="Family name"
                  placeholder="The Smiths"
                  value={createName}
                  onChange={(e) => setCreateName(e.target.value)}
                  required
                  fullWidth
                />
                <Button type="submit" variant="contained" size="large" disabled={onboardSubmitting}>
                  Create family
                </Button>
              </Stack>
            </Paper>
          </Grid>
          <Grid size={{ xs: 12, sm: 6 }}>
            <Paper sx={{ p: 3, height: "100%" }}>
              <Stack spacing={2} component="form" onSubmit={handleJoin} sx={{ height: "100%" }}>
                <Box sx={{ flexGrow: 1 }}>
                  <Typography variant="h6" component="h2">
                    Join a family
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Already have an invite code from a sibling?
                  </Typography>
                </Box>
                <TextField
                  label="Invite code"
                  value={joinCode}
                  onChange={(e) => setJoinCode(e.target.value)}
                  required
                  fullWidth
                />
                <Button type="submit" variant="outlined" size="large" disabled={onboardSubmitting}>
                  Join family
                </Button>
              </Stack>
            </Paper>
          </Grid>
        </Grid>
      </>
    );
  }

  const colors = familyColorMap(families);

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Welcome{user ? `, ${user.name}` : ""}
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Stack spacing={2}>
        {families.map((family) => (
          <FamilyDashboardBlock
            key={family.id}
            family={family}
            color={colors.get(family.id)!}
            stats={stats[family.id]}
            activity={activity[family.id]}
            contributions={contributions[family.id]}
          />
        ))}
      </Stack>
    </>
  );
}

function FamilyDashboardBlock({
  family,
  color,
  stats,
  activity,
  contributions,
}: {
  family: Family;
  color: FamilyColor;
  stats?: FamilyStats;
  activity?: ActivityEntry[];
  contributions?: Contribution[];
}) {
  // A "fresh" family (no parent tracked yet, no tasks ever created) gets a
  // short checklist pointing at the two things that actually unlock the
  // rest of the app, instead of three stat tiles that all read zero.
  const isFresh = stats && stats.parents === 0 && stats.openTasks === 0 && stats.documents === 0;

  return (
    <FamilySection
      family={family}
      color={color}
      action={
        <Stack direction="row" spacing={2} sx={{ alignItems: "center" }}>
          <Button size="small" startIcon={<AddIcon />} component={RouterLink} to="/tasks" sx={{ color: color.text }}>
            Task
          </Button>
          <Button size="small" startIcon={<AddIcon />} component={RouterLink} to="/expenses" sx={{ color: color.text }}>
            Expense
          </Button>
          <RouterLink to="/family" style={{ fontSize: 14, color: color.text, fontWeight: 500 }}>
            View family
          </RouterLink>
        </Stack>
      }
    >
      {isFresh ? (
        <List sx={{ py: 1 }}>
          <ChecklistItem
            done={false}
            icon={<PersonAddAlt1Icon fontSize="small" />}
            label="Add a parent to start tracking care for them"
            to="/parents"
          />
          <ChecklistItem
            done={false}
            icon={<GroupAddIcon fontSize="small" />}
            label="Invite a sibling to join this family"
            to="/family"
          />
          <ChecklistItem
            done={false}
            icon={<AssignmentIcon fontSize="small" />}
            label="Create your first task"
            to="/tasks"
          />
        </List>
      ) : (
        <>
        <Grid container spacing={0} sx={{ p: 2 }}>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
              <AssignmentIcon fontSize="small" sx={{ color: "text.disabled" }} />
              <Box>
                <Typography variant="overline" color="text.secondary" sx={{ lineHeight: 1 }}>
                  Open tasks
                </Typography>
                <Typography variant="h5">{stats ? stats.openTasks : "—"}</Typography>
              </Box>
            </Stack>
          </Grid>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
              <AttachMoneyIcon fontSize="small" sx={{ color: "text.disabled" }} />
              <Box>
                <Typography variant="overline" color="text.secondary" sx={{ lineHeight: 1 }}>
                  Your balance
                </Typography>
                <Typography
                  variant="h5"
                  sx={stats?.balanceCents ? { color: stats.balanceCents > 0 ? "success.main" : "error.main" } : undefined}
                >
                  {stats ? (stats.balanceCents ? formatCents(Math.abs(stats.balanceCents)) : "Settled") : "—"}
                </Typography>
              </Box>
            </Stack>
          </Grid>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
              <DescriptionIcon fontSize="small" sx={{ color: "text.disabled" }} />
              <Box>
                <Typography variant="overline" color="text.secondary" sx={{ lineHeight: 1 }}>
                  Documents
                </Typography>
                <Typography variant="h5">{stats ? stats.documents : "—"}</Typography>
              </Box>
            </Stack>
          </Grid>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
              <PersonAddAlt1Icon fontSize="small" sx={{ color: "text.disabled" }} />
              <Box>
                <Typography variant="overline" color="text.secondary" sx={{ lineHeight: 1 }}>
                  Parents
                </Typography>
                <Typography variant="h5">{stats ? stats.parents : "—"}</Typography>
              </Box>
            </Stack>
          </Grid>
        </Grid>
        <Divider />
        <Box sx={{ px: 1, pt: 1 }}>
          <Typography variant="overline" color="text.secondary" sx={{ pl: 1 }}>
            Recent activity
          </Typography>
          <ActivityFeed entries={activity ?? []} />
        </Box>
        <Divider />
        <Box sx={{ px: 1, py: 1 }}>
          <Typography variant="overline" color="text.secondary" sx={{ pl: 1 }}>
            Who's contributing
          </Typography>
          <ContributionLeaderboard entries={contributions ?? []} />
        </Box>
        </>
      )}
    </FamilySection>
  );
}

function ChecklistItem({ done, icon, label, to }: { done: boolean; icon: ReactNode; label: string; to: string }) {
  return (
    <ListItem
      component={RouterLink}
      to={to}
      sx={{ color: "inherit", textDecoration: "none", "&:hover": { bgcolor: "action.hover" } }}
    >
      <ListItemIcon sx={{ minWidth: 36, color: done ? "success.main" : "text.disabled" }}>
        {done ? <CheckCircleIcon fontSize="small" /> : <RadioButtonUncheckedIcon fontSize="small" />}
      </ListItemIcon>
      <ListItemText
        primary={label}
        slotProps={{ primary: { sx: { textDecoration: done ? "line-through" : "none" } } }}
      />
      {icon}
    </ListItem>
  );
}
