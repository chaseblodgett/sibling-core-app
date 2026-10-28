import AddIcon from "@mui/icons-material/Add";
import AssignmentIcon from "@mui/icons-material/Assignment";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import DownloadIcon from "@mui/icons-material/Download";
import RepeatIcon from "@mui/icons-material/Repeat";
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  IconButton,
  List,
  ListItem,
  ListItemText,
  MenuItem,
  Paper,
  Select,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { type FormEvent, useCallback, useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { ApiError } from "../api/client";
import { exportApi } from "../api/export";
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { tasksApi } from "../api/tasks";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { EmptyState } from "../components/EmptyState";
import { FamilySection } from "../components/FamilySection";
import { useAuth } from "../context/AuthContext";
import { useFamilies } from "../context/FamilyContext";
import { useNotify } from "../context/NotificationContext";
import { useConfirm } from "../hooks/useConfirm";
import { familyColorMap } from "../utils/familyColor";
import type { Family, Member, Parent, RecurrenceRule, Task, TaskStatus } from "../types";

const STATUS_LABEL: Record<TaskStatus, string> = {
  todo: "To do",
  in_progress: "In progress",
  done: "Done",
};

const RECURRENCE_LABEL: Record<RecurrenceRule, string> = {
  none: "Doesn't repeat",
  daily: "Repeats daily",
  weekly: "Repeats weekly",
  monthly: "Repeats monthly",
};

function todayISODate() {
  const d = new Date();
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
}

function isOverdue(dueDate: string) {
  return dueDate.slice(0, 10) < todayISODate();
}

function formatDueDate(dueDate: string) {
  return new Date(dueDate).toLocaleDateString(undefined, { timeZone: "UTC", month: "short", day: "numeric" });
}

export function TasksPage() {
  const { user } = useAuth();
  const { families, loading: familiesLoading } = useFamilies();
  const notify = useNotify();
  const { state: confirmState, confirm, close: closeConfirm, handleConfirm } = useConfirm();
  const [searchParams, setSearchParams] = useSearchParams();
  const parentFilter = searchParams.get("parent_id") ?? "";

  const [tasks, setTasks] = useState<Task[]>([]);
  const [parents, setParents] = useState<Parent[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [createFamilyId, setCreateFamilyId] = useState("");
  const [parentId, setParentId] = useState("");
  const [assigneeId, setAssigneeId] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [recurrenceRule, setRecurrenceRule] = useState<RecurrenceRule>("none");
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async (targets: Family[], parentIdFilter: string) => {
    if (targets.length === 0) {
      setTasks([]);
      setParents([]);
      setMembers([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const [taskLists, parentLists, memberLists] = await Promise.all([
        Promise.all(targets.map((f) => tasksApi.list(f.id, parentIdFilter || undefined))),
        Promise.all(targets.map((f) => parentsApi.list(f.id))),
        Promise.all(targets.map((f) => familyApi.listMembers(f.id))),
      ]);
      setTasks(taskLists.flat());
      setParents(parentLists.flat());
      setMembers(memberLists.flat());
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't load tasks.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!familiesLoading) {
      load(families, parentFilter);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [familiesLoading, families, parentFilter]);

  function isViewerOf(familyId: string) {
    return members.some((m) => m.user_id === user?.id && m.family_id === familyId && m.role === "viewer");
  }
  const writableFamilies = families.filter((f) => !isViewerOf(f.id));

  useEffect(() => {
    setCreateFamilyId((prev) => prev || writableFamilies[0]?.id || "");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [families, members]);

  const createFamilyParents = parents.filter((p) => p.family_id === createFamilyId);
  const createFamilyMembers = members.filter((m) => m.family_id === createFamilyId);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    if (!createFamilyId) return;
    if (recurrenceRule !== "none" && !dueDate) {
      setError("Repeating tasks need a due date — that's what the next occurrence is scheduled from.");
      return;
    }
    setError(null);
    setSubmitting(true);
    try {
      await tasksApi.create(createFamilyId, {
        title,
        description,
        parent_id: parentId || undefined,
        assignee_id: assigneeId || undefined,
        due_date: dueDate ? `${dueDate}T00:00:00Z` : undefined,
        recurrence_rule: recurrenceRule !== "none" ? recurrenceRule : undefined,
      });
      setTitle("");
      setDescription("");
      setParentId("");
      setAssigneeId("");
      setDueDate("");
      setRecurrenceRule("none");
      notify("Task created.");
      await load(families, parentFilter);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't create that task.");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleStatusChange(task: Task, status: TaskStatus) {
    setError(null);
    try {
      await tasksApi.update(task.id, { status });
      // Full reload rather than a local splice: completing a recurring
      // task spawns a new row server-side, which a local update can't see.
      await load(families, parentFilter);
      if (status === "done") {
        notify(
          task.recurrence_rule !== "none"
            ? "Nice work — task marked done. Next occurrence scheduled."
            : "Nice work — task marked done.",
        );
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't update that task.");
    }
  }

  function confirmDelete(task: Task) {
    confirm({
      title: "Delete task?",
      message: `Delete "${task.title}"? This can't be undone.`,
      confirmLabel: "Delete",
      destructive: true,
      onConfirm: () => handleDelete(task.id),
    });
  }

  async function handleDelete(taskId: string) {
    setError(null);
    try {
      await tasksApi.delete(taskId);
      setTasks((prev) => prev.filter((t) => t.id !== taskId));
      notify("Task deleted.");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't delete that task.");
    }
  }

  function clearParentFilter() {
    const next = new URLSearchParams(searchParams);
    next.delete("parent_id");
    setSearchParams(next);
  }

  function parentName(id?: string) {
    return parents.find((p) => p.id === id)?.name;
  }

  function memberName(id?: string) {
    const member = members.find((m) => m.id === id);
    return member?.full_name || member?.email;
  }

  if (familiesLoading || loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 8 }}>
        <CircularProgress sx={{ color: "primary.main" }} />
      </Box>
    );
  }

  if (families.length === 0) {
    return (
      <EmptyState
        icon={<AssignmentIcon />}
        title="No family yet"
        description="Create or join a family first, then come back here to start tracking tasks."
      />
    );
  }

  const colors = familyColorMap(families);
  const filteredParentName = parentFilter ? parentName(parentFilter) : undefined;

  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2, flexWrap: "wrap", gap: 2 }}>
        <Typography variant="h4" component="h1">
          Tasks
        </Typography>
        {parentFilter && (
          <Chip label={`Filtered by: ${filteredParentName ?? "parent"}`} onDelete={clearParentFilter} color="primary" />
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Stack spacing={2}>
        {writableFamilies.length > 0 && (
        <Paper sx={{ p: 3 }}>
          <Stack spacing={2} component="form" onSubmit={handleCreate}>
            <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap" }}>
              <TextField
                label="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                size="small"
                sx={{ flex: 1, minWidth: 200 }}
              />
              <TextField
                label="Description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                size="small"
                sx={{ flex: 1, minWidth: 200 }}
              />
            </Stack>
            <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap", alignItems: "center" }}>
              {writableFamilies.length > 1 && (
                <Select
                  size="small"
                  value={createFamilyId}
                  onChange={(e) => {
                    setCreateFamilyId(e.target.value);
                    setParentId("");
                    setAssigneeId("");
                  }}
                  sx={{ minWidth: 180 }}
                >
                  {writableFamilies.map((f) => (
                    <MenuItem key={f.id} value={f.id}>
                      {f.name}
                    </MenuItem>
                  ))}
                </Select>
              )}
              <Select
                size="small"
                value={parentId}
                onChange={(e) => setParentId(e.target.value)}
                displayEmpty
                sx={{ minWidth: 180 }}
              >
                <MenuItem value="">No parent (general)</MenuItem>
                {createFamilyParents.map((p) => (
                  <MenuItem key={p.id} value={p.id}>
                    {p.name}
                  </MenuItem>
                ))}
              </Select>
              <Select
                size="small"
                value={assigneeId}
                onChange={(e) => setAssigneeId(e.target.value)}
                displayEmpty
                sx={{ minWidth: 180 }}
              >
                <MenuItem value="">Unassigned</MenuItem>
                {createFamilyMembers.map((m) => (
                  <MenuItem key={m.id} value={m.id}>
                    {m.full_name || m.email}
                  </MenuItem>
                ))}
              </Select>
              <TextField
                label="Due date"
                type="date"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
                size="small"
                sx={{ minWidth: 160 }}
                slotProps={{ inputLabel: { shrink: true } }}
              />
              <Select
                size="small"
                value={recurrenceRule}
                onChange={(e) => setRecurrenceRule(e.target.value as RecurrenceRule)}
                sx={{ minWidth: 180 }}
              >
                {(Object.keys(RECURRENCE_LABEL) as RecurrenceRule[]).map((rule) => (
                  <MenuItem key={rule} value={rule}>
                    {RECURRENCE_LABEL[rule]}
                  </MenuItem>
                ))}
              </Select>
              <Button type="submit" variant="contained" startIcon={<AddIcon />} disabled={submitting}>
                New task
              </Button>
            </Stack>
          </Stack>
        </Paper>
        )}

        {families
          .filter((f) => !parentFilter || tasks.some((t) => t.family_id === f.id))
          .map((family) => {
            // Open tasks first, then done ones at the bottom; within each
            // group, soonest-due first and undated tasks last.
            const familyTasks = [...tasks]
              .filter((t) => t.family_id === family.id)
              .sort((a, b) => {
                const doneDiff = Number(a.status === "done") - Number(b.status === "done");
                if (doneDiff !== 0) return doneDiff;
                if (a.due_date && b.due_date) return a.due_date < b.due_date ? -1 : a.due_date > b.due_date ? 1 : 0;
                if (a.due_date) return -1;
                if (b.due_date) return 1;
                return 0;
              });
            return (
              <FamilySection
                key={family.id}
                family={family}
                color={colors.get(family.id)!}
                action={
                  <Button
                    size="small"
                    startIcon={<DownloadIcon />}
                    onClick={() => exportApi.tasks(family.id).catch(() => notify("Couldn't export tasks.", "error"))}
                    sx={{ color: colors.get(family.id)!.text }}
                  >
                    Export CSV
                  </Button>
                }
              >
                {familyTasks.length === 0 ? (
                  <EmptyState
                    icon={<AssignmentIcon />}
                    title="No tasks yet"
                    description="Tasks assigned to family members will show up here — add one using the form above."
                  />
                ) : (
                  <List>
                    {familyTasks.map((task) => (
                      <ListItem
                        key={task.id}
                        secondaryAction={
                          !isViewerOf(task.family_id) && (
                            <IconButton edge="end" aria-label="Delete task" onClick={() => confirmDelete(task)}>
                              <DeleteOutlineIcon fontSize="small" />
                            </IconButton>
                          )
                        }
                      >
                        <ListItemText
                          primary={
                            <Stack direction="row" spacing={1} sx={{ alignItems: "center", flexWrap: "wrap" }}>
                              <span style={{ textDecoration: task.status === "done" ? "line-through" : "none" }}>
                                {task.title}
                              </span>
                              {parentName(task.parent_id) && <Chip label={parentName(task.parent_id)} size="small" />}
                              {memberName(task.assignee_id) && (
                                <Chip label={memberName(task.assignee_id)} size="small" variant="outlined" />
                              )}
                              {task.due_date && (
                                <Chip
                                  label={formatDueDate(task.due_date)}
                                  size="small"
                                  color={task.status !== "done" && isOverdue(task.due_date) ? "error" : "default"}
                                  variant="outlined"
                                />
                              )}
                              {task.recurrence_rule !== "none" && (
                                <Chip
                                  icon={<RepeatIcon fontSize="small" />}
                                  label={task.recurrence_rule}
                                  size="small"
                                  variant="outlined"
                                />
                              )}
                            </Stack>
                          }
                          secondary={task.description || undefined}
                        />
                        <Select
                          size="small"
                          value={task.status}
                          onChange={(e) => handleStatusChange(task, e.target.value as TaskStatus)}
                          disabled={isViewerOf(task.family_id)}
                          sx={{ minWidth: 140, mr: 2 }}
                        >
                          {Object.entries(STATUS_LABEL).map(([value, label]) => (
                            <MenuItem key={value} value={value}>
                              {label}
                            </MenuItem>
                          ))}
                        </Select>
                      </ListItem>
                    ))}
                  </List>
                )}
              </FamilySection>
            );
          })}
      </Stack>

      {confirmState && (
        <ConfirmDialog
          open
          title={confirmState.title}
          message={confirmState.message}
          confirmLabel={confirmState.confirmLabel}
          destructive={confirmState.destructive}
          onConfirm={handleConfirm}
          onCancel={closeConfirm}
        />
      )}
    </>
  );
}
