import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  IconButton,
  Link,
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
import { Link as RouterLink } from "react-router-dom";
import { ApiError } from "../api/client";
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { useAuth } from "../context/AuthContext";
import { useFamilies } from "../context/FamilyContext";
import type { Member, Parent } from "../types";

export function ParentsPage() {
  const { user } = useAuth();
  const { families, loading: familiesLoading, selectedFamilyId } = useFamilies();

  const [parents, setParents] = useState<Parent[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [name, setName] = useState("");
  const [relation, setRelation] = useState("");
  const [createFamilyId, setCreateFamilyId] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const viewingAllFamilies = selectedFamilyId === null;
  const familiesToLoad = viewingAllFamilies ? families : families.filter((f) => f.id === selectedFamilyId);

  const load = useCallback(async (targets: typeof families) => {
    if (targets.length === 0) {
      setParents([]);
      setMembers([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const [parentLists, memberLists] = await Promise.all([
        Promise.all(targets.map((f) => parentsApi.list(f.id))),
        Promise.all(targets.map((f) => familyApi.listMembers(f.id))),
      ]);
      setParents(parentLists.flat());
      setMembers(memberLists.flat());
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't load parents.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!familiesLoading) {
      load(familiesToLoad);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [familiesLoading, selectedFamilyId, families]);

  useEffect(() => {
    setCreateFamilyId(selectedFamilyId ?? families[0]?.id ?? "");
  }, [selectedFamilyId, families]);

  const isOwnerOf = useCallback(
    (familyId: string) => members.some((m) => m.user_id === user?.id && m.family_id === familyId && m.role === "owner"),
    [members, user?.id],
  );

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    if (!createFamilyId) return;
    setError(null);
    setSubmitting(true);
    try {
      await parentsApi.create(createFamilyId, { name, relation });
      setName("");
      setRelation("");
      await load(familiesToLoad);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't add that parent.");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(parentId: string) {
    setError(null);
    try {
      await parentsApi.delete(parentId);
      setParents((prev) => prev.filter((p) => p.id !== parentId));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't remove that parent.");
    }
  }

  function familyName(id: string) {
    return families.find((f) => f.id === id)?.name;
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
      <Alert severity="info">
        You need a family first — <Link component={RouterLink} to="/family">create or join one</Link> to start
        tracking parents.
      </Alert>
    );
  }

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Parents
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        The people your family is coordinating care for. Tasks, expenses, and documents can be tied to a parent.
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Stack spacing={2}>
        <Paper sx={{ p: 3 }}>
          <Stack direction="row" spacing={2} component="form" onSubmit={handleCreate} sx={{ flexWrap: "wrap" }}>
            {viewingAllFamilies && families.length > 1 && (
              <Select size="small" value={createFamilyId} onChange={(e) => setCreateFamilyId(e.target.value)}>
                {families.map((f) => (
                  <MenuItem key={f.id} value={f.id}>
                    {f.name}
                  </MenuItem>
                ))}
              </Select>
            )}
            <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} required size="small" />
            <TextField
              label="Relation"
              placeholder="Mom, Dad, Stepmom…"
              value={relation}
              onChange={(e) => setRelation(e.target.value)}
              size="small"
            />
            <Button type="submit" variant="contained" disabled={submitting}>
              Add parent
            </Button>
          </Stack>
        </Paper>

        <Paper>
          <List>
            {parents.length === 0 && (
              <ListItem>
                <ListItemText primary="No parents yet" secondary="Add one above to start tracking tasks for them." />
              </ListItem>
            )}
            {parents.map((parent) => (
              <ListItem
                key={parent.id}
                secondaryAction={
                  isOwnerOf(parent.family_id) && (
                    <IconButton edge="end" aria-label="Remove parent" onClick={() => handleDelete(parent.id)}>
                      <DeleteOutlineIcon fontSize="small" />
                    </IconButton>
                  )
                }
              >
                <ListItemText
                  primary={
                    <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
                      <span>{parent.name}</span>
                      {parent.relation && <Chip label={parent.relation} size="small" />}
                      {viewingAllFamilies && (
                        <Chip label={familyName(parent.family_id)} size="small" color="primary" variant="outlined" />
                      )}
                    </Stack>
                  }
                  secondary={
                    <Link component={RouterLink} to={`/tasks?parent_id=${parent.id}&family_id=${parent.family_id}`}>
                      View tasks
                    </Link>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Paper>
      </Stack>
    </>
  );
}
