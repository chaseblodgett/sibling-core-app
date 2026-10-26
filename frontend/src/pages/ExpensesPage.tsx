import AddIcon from "@mui/icons-material/Add";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import {
  Alert,
  Box,
  Button,
  Checkbox,
  Chip,
  CircularProgress,
  Collapse,
  FormControlLabel,
  Grid,
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
import { ApiError } from "../api/client";
import { expensesApi } from "../api/expenses";
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { useAuth } from "../context/AuthContext";
import { useFamilies } from "../context/FamilyContext";
import type { Expense, ExpenseSplit, Member, Parent } from "../types";

function formatCents(cents: number, currency = "USD") {
  return (cents / 100).toLocaleString(undefined, { style: "currency", currency });
}

export function ExpensesPage() {
  const { user } = useAuth();
  const { families, loading: familiesLoading, selectedFamilyId } = useFamilies();

  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [parents, setParents] = useState<Parent[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [balances, setBalances] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [description, setDescription] = useState("");
  const [amountDollars, setAmountDollars] = useState("");
  const [category, setCategory] = useState("");
  const [createFamilyId, setCreateFamilyId] = useState("");
  const [parentId, setParentId] = useState("");
  const [paidBy, setPaidBy] = useState("");
  const [splitMemberIds, setSplitMemberIds] = useState<string[]>([]);
  const [submitting, setSubmitting] = useState(false);

  const viewingAllFamilies = selectedFamilyId === null;
  const familiesToLoad = viewingAllFamilies ? families : families.filter((f) => f.id === selectedFamilyId);

  const load = useCallback(async (targets: typeof families, allMode: boolean) => {
    if (targets.length === 0) {
      setExpenses([]);
      setParents([]);
      setMembers([]);
      setBalances({});
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const [expenseLists, parentLists, memberLists] = await Promise.all([
        Promise.all(targets.map((f) => expensesApi.list(f.id))),
        Promise.all(targets.map((f) => parentsApi.list(f.id))),
        Promise.all(targets.map((f) => familyApi.listMembers(f.id))),
      ]);
      setExpenses(expenseLists.flat());
      setParents(parentLists.flat());
      setMembers(memberLists.flat());

      if (!allMode && targets.length === 1) {
        const balanceList = await expensesApi.balances(targets[0].id);
        setBalances(Object.fromEntries(balanceList.map((b) => [b.member_id, b.net_cents])));
      } else {
        setBalances({});
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't load expenses.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!familiesLoading) {
      load(familiesToLoad, viewingAllFamilies);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [familiesLoading, selectedFamilyId, families]);

  useEffect(() => {
    setCreateFamilyId(selectedFamilyId ?? families[0]?.id ?? "");
  }, [selectedFamilyId, families]);

  const createFamilyParents = parents.filter((p) => p.family_id === createFamilyId);
  const createFamilyMembers = members.filter((m) => m.family_id === createFamilyId);

  useEffect(() => {
    // Default: split evenly among everyone in the chosen family, and pay
    // as yourself — the common case for logging your own expense.
    setSplitMemberIds(createFamilyMembers.map((m) => m.id));
    setPaidBy(createFamilyMembers.find((m) => m.user_id === user?.id)?.id ?? "");
  }, [createFamilyId, createFamilyMembers, user?.id]);

  function toggleSplitMember(memberId: string) {
    setSplitMemberIds((prev) => (prev.includes(memberId) ? prev.filter((id) => id !== memberId) : [...prev, memberId]));
  }

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    if (!createFamilyId) return;
    const amountCents = Math.round(parseFloat(amountDollars) * 100);
    if (!Number.isFinite(amountCents) || amountCents <= 0) {
      setError("Enter a valid amount greater than zero.");
      return;
    }

    setError(null);
    setSubmitting(true);
    try {
      const splits =
        splitMemberIds.length === 0
          ? []
          : (() => {
              const share = Math.floor(amountCents / splitMemberIds.length);
              const remainder = amountCents - share * splitMemberIds.length;
              return splitMemberIds.map((memberId, i) => ({
                member_id: memberId,
                share_cents: i === 0 ? share + remainder : share,
              }));
            })();

      await expensesApi.create(createFamilyId, {
        description,
        amount_cents: amountCents,
        currency: "USD",
        category,
        parent_id: parentId || undefined,
        paid_by: paidBy || undefined,
        splits,
      });
      setDescription("");
      setAmountDollars("");
      setCategory("");
      setParentId("");
      await load(familiesToLoad, viewingAllFamilies);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't log that expense.");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(expenseId: string) {
    setError(null);
    try {
      await expensesApi.delete(expenseId);
      setExpenses((prev) => prev.filter((e) => e.id !== expenseId));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't delete that expense.");
    }
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
    return <Alert severity="info">You need a family first to track expenses.</Alert>;
  }

  const hasBalances = Object.keys(balances).length > 0;

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Expenses
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Grid container spacing={2}>
        <Grid size={{ xs: 12, md: 8 }}>
          <Stack spacing={2}>
            <Paper sx={{ p: 3 }}>
              <Stack spacing={2} component="form" onSubmit={handleCreate}>
                <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap" }}>
                  <TextField
                    label="Description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    required
                    size="small"
                    sx={{ flex: 1, minWidth: 180 }}
                  />
                  <TextField
                    label="Amount"
                    type="number"
                    slotProps={{ htmlInput: { step: "0.01", min: "0" } }}
                    value={amountDollars}
                    onChange={(e) => setAmountDollars(e.target.value)}
                    required
                    size="small"
                    sx={{ width: 140 }}
                  />
                  <TextField
                    label="Category"
                    value={category}
                    onChange={(e) => setCategory(e.target.value)}
                    size="small"
                    sx={{ width: 160 }}
                  />
                </Stack>
                <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap", alignItems: "center" }}>
                  {viewingAllFamilies && families.length > 1 && (
                    <Select size="small" value={createFamilyId} onChange={(e) => setCreateFamilyId(e.target.value)}>
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
                    sx={{ minWidth: 160 }}
                  >
                    <MenuItem value="">No parent (general)</MenuItem>
                    {createFamilyParents.map((p) => (
                      <MenuItem key={p.id} value={p.id}>
                        {p.name}
                      </MenuItem>
                    ))}
                  </Select>
                  <Select size="small" value={paidBy} onChange={(e) => setPaidBy(e.target.value)} sx={{ minWidth: 160 }}>
                    <MenuItem value="" disabled>
                      Paid by…
                    </MenuItem>
                    {createFamilyMembers.map((m) => (
                      <MenuItem key={m.id} value={m.id}>
                        {m.full_name || m.email}
                      </MenuItem>
                    ))}
                  </Select>
                </Stack>
                <Box>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                    Split evenly between:
                  </Typography>
                  <Stack direction="row" sx={{ flexWrap: "wrap" }}>
                    {createFamilyMembers.map((m) => (
                      <FormControlLabel
                        key={m.id}
                        control={
                          <Checkbox
                            size="small"
                            checked={splitMemberIds.includes(m.id)}
                            onChange={() => toggleSplitMember(m.id)}
                          />
                        }
                        label={m.full_name || m.email}
                      />
                    ))}
                  </Stack>
                </Box>
                <Box>
                  <Button type="submit" variant="contained" startIcon={<AddIcon />} disabled={submitting}>
                    Log expense
                  </Button>
                </Box>
              </Stack>
            </Paper>

            <Paper>
              <List>
                {expenses.length === 0 && (
                  <ListItem>
                    <ListItemText primary="No expenses yet" secondary="Logged expenses and their splits will show up here." />
                  </ListItem>
                )}
                {expenses.map((expense) => (
                  <ExpenseRow
                    key={expense.id}
                    expense={expense}
                    familyLabel={viewingAllFamilies ? familyName(expense.family_id) : undefined}
                    parentLabel={parentName(expense.parent_id)}
                    paidByLabel={memberName(expense.paid_by)}
                    memberName={memberName}
                    onDelete={() => handleDelete(expense.id)}
                  />
                ))}
              </List>
            </Paper>
          </Stack>
        </Grid>

        <Grid size={{ xs: 12, md: 4 }}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" component="h2" gutterBottom>
              Balances
            </Typography>
            {viewingAllFamilies ? (
              <Typography variant="body2" color="text.secondary">
                Pick a specific family above to see its balances.
              </Typography>
            ) : !hasBalances ? (
              <Typography variant="body2" color="text.secondary">
                All settled up.
              </Typography>
            ) : (
              <List dense>
                {Object.entries(balances).map(([memberId, net]) => (
                  <ListItem key={memberId} disableGutters>
                    <ListItemText
                      primary={memberName(memberId)}
                      secondary={net > 0 ? "is owed" : "owes"}
                    />
                    <Typography variant="body2" sx={{ color: net > 0 ? "success.main" : "error.main", fontWeight: 600 }}>
                      {formatCents(Math.abs(net))}
                    </Typography>
                  </ListItem>
                ))}
              </List>
            )}
          </Paper>
        </Grid>
      </Grid>
    </>
  );
}

function ExpenseRow({
  expense,
  familyLabel,
  parentLabel,
  paidByLabel,
  memberName,
  onDelete,
}: {
  expense: Expense;
  familyLabel?: string;
  parentLabel?: string;
  paidByLabel?: string;
  memberName: (id?: string) => string | undefined;
  onDelete: () => void;
}) {
  const [expanded, setExpanded] = useState(false);
  const [splits, setSplits] = useState<ExpenseSplit[]>([]);
  const [splitsLoaded, setSplitsLoaded] = useState(false);

  async function toggleExpand() {
    setExpanded((prev) => !prev);
    if (!splitsLoaded) {
      try {
        setSplits(await expensesApi.listSplits(expense.id));
        setSplitsLoaded(true);
      } catch {
        // best-effort — the row still works without split detail
      }
    }
  }

  async function handleSettle(splitId: string) {
    const updated = await expensesApi.settleSplit(splitId);
    setSplits((prev) => prev.map((s) => (s.id === splitId ? updated : s)));
  }

  return (
    <>
      <ListItem
        secondaryAction={
          <Stack direction="row" spacing={0.5}>
            <IconButton edge="end" aria-label="Show splits" onClick={toggleExpand}>
              <ExpandMoreIcon
                fontSize="small"
                sx={{ transform: expanded ? "rotate(180deg)" : "none", transition: "transform 0.15s" }}
              />
            </IconButton>
            <IconButton edge="end" aria-label="Delete expense" onClick={onDelete}>
              <DeleteOutlineIcon fontSize="small" />
            </IconButton>
          </Stack>
        }
      >
        <ListItemText
          primary={
            <Stack direction="row" spacing={1} sx={{ alignItems: "center", flexWrap: "wrap" }}>
              <span>{expense.description}</span>
              <Chip label={formatCents(expense.amount_cents, expense.currency)} size="small" color="primary" />
              {expense.category && <Chip label={expense.category} size="small" variant="outlined" />}
              {familyLabel && <Chip label={familyLabel} size="small" color="primary" variant="outlined" />}
              {parentLabel && <Chip label={parentLabel} size="small" />}
            </Stack>
          }
          secondary={paidByLabel ? `Paid by ${paidByLabel}` : undefined}
        />
      </ListItem>
      <Collapse in={expanded} unmountOnExit>
        <List dense sx={{ pl: 4, pr: 2 }}>
          {splits.map((split) => (
            <ListItem key={split.id} disableGutters>
              <ListItemText primary={memberName(split.member_id)} secondary={`Owes ${formatCents(split.share_cents)}`} />
              {split.settled ? (
                <Chip label="Settled" size="small" color="success" variant="outlined" />
              ) : (
                <Button size="small" onClick={() => handleSettle(split.id)}>
                  Mark settled
                </Button>
              )}
            </ListItem>
          ))}
          {splitsLoaded && splits.length === 0 && (
            <ListItem disableGutters>
              <ListItemText secondary="No splits — this expense isn't owed by anyone." />
            </ListItem>
          )}
        </List>
      </Collapse>
    </>
  );
}
