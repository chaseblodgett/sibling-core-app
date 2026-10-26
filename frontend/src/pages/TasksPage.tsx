import AddIcon from "@mui/icons-material/Add";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
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
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { tasksApi } from "../api/tasks";
import { useFamilies } from "../context/FamilyContext";
import type { Member, Parent, Task, TaskStatus } from "../types";

const STATUS_LABEL: Record<TaskStatus, string> = {
  todo: "To do",
  in_progress: "In progress",
  done: "Done",
};

export function TasksPage() {
  const { families, loading: familiesLoading, selectedFamilyId, setSelectedFamilyId } = useFamilies();
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
  const [submitting, setSubmitting] = useState(false);

  // A deep link from the Parents page ("view tasks" for a specific
  // parent) pins the family context so the parent filter actually matches
  // something, then falls back to the normal flow.
  useEffect(() => {
    const familyIdParam = searchParams.get("family_id");
    if (familyIdParam && familyIdParam !== selectedFamilyId) {
      setSelectedFamilyId(familyIdParam);
      const next = new URLSearchParams(searchParams);
      next.delete("family_id");
      setSearchParams(next, { replace: true });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  const viewingAllFamilies = selectedFamilyId === null;
  const familiesToLoad = viewingAllFamilies ? families : families.filter((f) => f.id === selectedFamilyId);

  const load = useCallback(async (targets: typeof families, parentIdFilter: string, allMode: boolean) => {
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
        Promise.all(targets.map((f) => tasksApi.list(f.id, allMode ? undefined : parentIdFilter || undefined))),
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
      load(familiesToLoad, parentFilter, viewingAllFamilies);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [familiesLoading, selectedFamilyId, families, parentFilter]);

  useEffect(() => {
    // Keep the create-form's target family sane as context changes: pin
    // to the selected family, or default to the first when viewing all.
    setCreateFamilyId(selectedFamilyId ?? families[0]?.id ?? "");
  }, [selectedFamilyId, families]);

  const createFamilyParents = parents.filter((p) => p.family_id === createFamilyId);
  const createFamilyMembers = members.filter((m) => m.family_id === createFamilyId);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    if (!createFamilyId) return;
    setError(null);
    setSubmitting(true);
    try {
      await tasksApi.create(createFamilyId, {
        title,
        description,
        parent_id: parentId || undefined,
        assignee_id: assigneeId || undefined,
      });
      setTitle("");
      setDescription("");
      setParentId("");
      setAssigneeId("");
      await load(familiesToLoad, parentFilter, viewingAllFamilies);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't create that task.");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleStatusChange(task: Task, status: TaskStatus) {
    setError(null);
    try {
      const updated = await tasksApi.update(task.id, { status });
      setTasks((prev) => prev.map((t) => (t.id === task.id ? updated : t)));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't update that task.");
    }
  }

  async function handleDelete(taskId: string) {
    setError(null);
    try {
      await tasksApi.delete(taskId);
      setTasks((prev) => prev.filter((t) => t.id !== taskId));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't delete that task.");
    }
  }

  function handleParentFilterChange(value: string) {
    const next = new URLSearchParams(searchParams);
    if (value) {
      next.set("parent_id", value);
    } else {
      next.delete("parent_id");
    }
    setSearchParams(next);
  }

  function familyName(id: string) {
    return families.find((f) => f.id === id)?.name;
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
    return <Alert severity="info">You need a family first to track tasks.</Alert>;
  }

  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2, flexWrap: "wrap", gap: 2 }}>
        <Typography variant="h4" component="h1">
          Tasks
        </Typography>
        {!viewingAllFamilies && (
          <Select
            size="small"
            value={parentFilter}
            onChange={(e) => handleParentFilterChange(e.target.value)}
            displayEmpty
            sx={{ minWidth: 200 }}
          >
            <MenuItem value="">All parents</MenuItem>
            {parents.map((p) => (
              <MenuItem key={p.id} value={p.id}>
                {p.name}
              </MenuItem>
            ))}
          </Select>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Stack spacing={2}>
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
              {viewingAllFamilies && families.length > 1 && (
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
                  {families.map((f) => (
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
              <Button type="submit" variant="contained" startIcon={<AddIcon />} disabled={submitting}>
                New task
              </Button>
            </Stack>
          </Stack>
        </Paper>

        <Paper>
          <List>
            {tasks.length === 0 && (
              <ListItem>
                <ListItemText primary="No tasks yet" secondary="Tasks assigned to family members will show up here." />
              </ListItem>
            )}
            {tasks.map((task) => (
              <ListItem
                key={task.id}
                secondaryAction={
                  <IconButton edge="end" aria-label="Delete task" onClick={() => handleDelete(task.id)}>
                    <DeleteOutlineIcon fontSize="small" />
                  </IconButton>
                }
              >
                <ListItemText
                  primary={
                    <Stack direction="row" spacing={1} sx={{ alignItems: "center", flexWrap: "wrap" }}>
                      <span>{task.title}</span>
                      {viewingAllFamilies && (
                        <Chip label={familyName(task.family_id)} size="small" color="primary" variant="outlined" />
                      )}
                      {parentName(task.parent_id) && <Chip label={parentName(task.parent_id)} size="small" />}
                      {memberName(task.assignee_id) && (
                        <Chip label={memberName(task.assignee_id)} size="small" variant="outlined" />
                      )}
                    </Stack>
                  }
                  secondary={task.description || undefined}
                />
                <Select
                  size="small"
                  value={task.status}
                  onChange={(e) => handleStatusChange(task, e.target.value as TaskStatus)}
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
        </Paper>
      </Stack>
    </>
  );
}
