import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import PersonAddAlt1Icon from "@mui/icons-material/PersonAddAlt1";
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
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { ApiError } from "../api/client";
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { EmptyState } from "../components/EmptyState";
import { FamilySection } from "../components/FamilySection";
import { useAuth } from "../context/AuthContext";
import { useFamilies } from "../context/FamilyContext";
import { useNotify } from "../context/NotificationContext";
import { useConfirm } from "../hooks/useConfirm";
import { familyColorMap } from "../utils/familyColor";
import type { Family, Member, Parent } from "../types";

export function ParentsPage() {
  const { user } = useAuth();
  const { families, loading: familiesLoading } = useFamilies();
  const notify = useNotify();
  const navigate = useNavigate();
  const { state: confirmState, confirm, close: closeConfirm, handleConfirm } = useConfirm();

  const [parents, setParents] = useState<Parent[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [name, setName] = useState("");
  const [relation, setRelation] = useState("");
  const [createFamilyId, setCreateFamilyId] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const load = useCallback(async (targets: Family[]) => {
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
      load(families);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [familiesLoading, families]);

  useEffect(() => {
    setCreateFamilyId((prev) => prev || families[0]?.id || "");
  }, [families]);

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
      notify(`${name} added.`);
      setName("");
      setRelation("");
      await load(families);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't add that parent.");
    } finally {
      setSubmitting(false);
    }
  }

  function confirmDelete(parent: Parent) {
    confirm({
      title: "Remove parent?",
      message: `Remove ${parent.name}? Tasks, expenses, and documents tied to them will become general (not deleted).`,
      confirmLabel: "Remove",
      destructive: true,
      onConfirm: () => handleDelete(parent.id),
    });
  }

  async function handleDelete(parentId: string) {
    setError(null);
    try {
      await parentsApi.delete(parentId);
      setParents((prev) => prev.filter((p) => p.id !== parentId));
      notify("Parent removed.");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't remove that parent.");
    }
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
        icon={<PersonAddAlt1Icon />}
        title="No family yet"
        description="Create or join a family first, then come back here to start tracking parents."
        actionLabel="Go to Family"
        onAction={() => navigate("/family")}
      />
    );
  }

  const colors = familyColorMap(families);

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
            {families.length > 1 && (
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

        {families.map((family) => {
          const familyParents = parents.filter((p) => p.family_id === family.id);
          return (
            <FamilySection key={family.id} family={family} color={colors.get(family.id)!}>
              {familyParents.length === 0 ? (
                <EmptyState
                  icon={<PersonAddAlt1Icon />}
                  title="No parents yet"
                  description="Add the person this family is coordinating care for, using the form above."
                />
              ) : (
                <List>
                  {familyParents.map((parent) => (
                    <ListItem
                      key={parent.id}
                      secondaryAction={
                        isOwnerOf(parent.family_id) && (
                          <IconButton edge="end" aria-label="Remove parent" onClick={() => confirmDelete(parent)}>
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
                          </Stack>
                        }
                        secondary={
                          <Link component={RouterLink} to={`/tasks?parent_id=${parent.id}`}>
                            View tasks
                          </Link>
                        }
                      />
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
