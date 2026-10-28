import AddIcon from "@mui/icons-material/Add";
import AttachFileIcon from "@mui/icons-material/AttachFile";
import AttachMoneyIcon from "@mui/icons-material/AttachMoney";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import DownloadIcon from "@mui/icons-material/Download";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import {
  Alert,
  Box,
  Button,
  Checkbox,
  Chip,
  CircularProgress,
  Collapse,
  Divider,
  FormControlLabel,
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
import { type ChangeEvent, type FormEvent, useCallback, useEffect, useRef, useState } from "react";
import { ApiError } from "../api/client";
import { documentsApi } from "../api/documents";
import { expensesApi } from "../api/expenses";
import { exportApi } from "../api/export";
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
import type { Balance, Document, Expense, ExpenseSplit, Family, Member, Parent } from "../types";

function formatCents(cents: number, currency = "USD") {
  return (cents / 100).toLocaleString(undefined, { style: "currency", currency });
}

export function ExpensesPage() {
  const { user } = useAuth();
  const { families, loading: familiesLoading } = useFamilies();
  const notify = useNotify();
  const { state: confirmState, confirm, close: closeConfirm, handleConfirm } = useConfirm();

  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [parents, setParents] = useState<Parent[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [balancesByFamily, setBalancesByFamily] = useState<Record<string, Balance[]>>({});
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

  const load = useCallback(async (targets: Family[]) => {
    if (targets.length === 0) {
      setExpenses([]);
      setParents([]);
      setMembers([]);
      setBalancesByFamily({});
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const [expenseLists, parentLists, memberLists, balanceLists] = await Promise.all([
        Promise.all(targets.map((f) => expensesApi.list(f.id))),
        Promise.all(targets.map((f) => parentsApi.list(f.id))),
        Promise.all(targets.map((f) => familyApi.listMembers(f.id))),
        Promise.all(targets.map((f) => expensesApi.balances(f.id))),
      ]);
      setExpenses(expenseLists.flat());
      setParents(parentLists.flat());
      setMembers(memberLists.flat());
      setBalancesByFamily(Object.fromEntries(targets.map((f, i) => [f.id, balanceLists[i]])));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't load expenses.");
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

  useEffect(() => {
    // Default: split evenly among everyone in the chosen family, and pay
    // as yourself — the common case for logging your own expense.
    // Depend on `members` (stable until actually reloaded) and filter
    // inside the effect — depending on `createFamilyMembers` directly
    // would re-run every render, since .filter() returns a new array
    // reference each time, causing an infinite render loop.
    const nextMembers = members.filter((m) => m.family_id === createFamilyId);
    setSplitMemberIds(nextMembers.map((m) => m.id));
    setPaidBy(nextMembers.find((m) => m.user_id === user?.id)?.id ?? "");
  }, [createFamilyId, members, user?.id]);

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
      notify("Expense logged.");
      await load(families);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't log that expense.");
    } finally {
      setSubmitting(false);
    }
  }

  function confirmDelete(expense: Expense) {
    confirm({
      title: "Delete expense?",
      message: `Delete "${expense.description}"? Any splits on it go with it. This can't be undone.`,
      confirmLabel: "Delete",
      destructive: true,
      onConfirm: () => handleDelete(expense.id),
    });
  }

  async function handleDelete(expenseId: string) {
    setError(null);
    try {
      await expensesApi.delete(expenseId);
      setExpenses((prev) => prev.filter((e) => e.id !== expenseId));
      notify("Expense deleted.");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't delete that expense.");
    }
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
        icon={<AttachMoneyIcon />}
        title="No family yet"
        description="Create or join a family first, then come back here to start tracking expenses."
      />
    );
  }

  const colors = familyColorMap(families);

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

      <Stack spacing={2}>
        {writableFamilies.length > 0 && (
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
              {writableFamilies.length > 1 && (
                <Select size="small" value={createFamilyId} onChange={(e) => setCreateFamilyId(e.target.value)}>
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
        )}

        {families.map((family) => {
          const familyExpenses = expenses.filter((e) => e.family_id === family.id);
          const familyBalances = balancesByFamily[family.id] ?? [];
          return (
            <FamilySection
              key={family.id}
              family={family}
              color={colors.get(family.id)!}
              action={
                <Button
                  size="small"
                  startIcon={<DownloadIcon />}
                  onClick={() => exportApi.expenses(family.id).catch(() => notify("Couldn't export expenses.", "error"))}
                  sx={{ color: colors.get(family.id)!.text }}
                >
                  Export CSV
                </Button>
              }
            >
              {familyBalances.length > 0 && (
                <>
                  <List dense sx={{ px: 1 }}>
                    {familyBalances.map((b) => (
                      <ListItem key={b.member_id} disableGutters sx={{ px: 2 }}>
                        <ListItemText primary={memberName(b.member_id)} secondary={b.net_cents > 0 ? "is owed" : "owes"} />
                        <Typography
                          variant="body2"
                          sx={{ color: b.net_cents > 0 ? "success.main" : "error.main", fontWeight: 600 }}
                        >
                          {formatCents(Math.abs(b.net_cents))}
                        </Typography>
                      </ListItem>
                    ))}
                  </List>
                  <Divider />
                </>
              )}
              {familyExpenses.length === 0 ? (
                <EmptyState
                  icon={<AttachMoneyIcon />}
                  title="No expenses yet"
                  description="Logged expenses and their splits will show up here — add one using the form above."
                />
              ) : (
                <List>
                  {familyExpenses.map((expense) => (
                    <ExpenseRow
                      key={expense.id}
                      expense={expense}
                      familyId={family.id}
                      isViewer={isViewerOf(family.id)}
                      parentLabel={parentName(expense.parent_id)}
                      paidByLabel={memberName(expense.paid_by)}
                      memberName={memberName}
                      onDelete={() => confirmDelete(expense)}
                    />
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

function ExpenseRow({
  expense,
  familyId,
  isViewer,
  parentLabel,
  paidByLabel,
  memberName,
  onDelete,
}: {
  expense: Expense;
  familyId: string;
  isViewer: boolean;
  parentLabel?: string;
  paidByLabel?: string;
  memberName: (id?: string) => string | undefined;
  onDelete: () => void;
}) {
  const notify = useNotify();
  const [expanded, setExpanded] = useState(false);
  const [splits, setSplits] = useState<ExpenseSplit[]>([]);
  const [splitsLoaded, setSplitsLoaded] = useState(false);

  const [documentId, setDocumentId] = useState(expense.document_id);
  const [receipt, setReceipt] = useState<Document | null>(null);
  const [attaching, setAttaching] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!documentId) {
      setReceipt(null);
      return;
    }
    let cancelled = false;
    documentsApi
      .get(documentId)
      .then((doc) => {
        if (!cancelled) setReceipt(doc);
      })
      .catch(() => {
        // best-effort — the row still works without the receipt's filename
      });
    return () => {
      cancelled = true;
    };
  }, [documentId]);

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
    notify("Marked settled.");
  }

  async function handleReceiptChosen(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    e.target.value = "";
    if (!file) return;

    setAttaching(true);
    try {
      const uploaded = await documentsApi.upload(familyId, file, expense.parent_id);
      const updated = await expensesApi.attachReceipt(expense.id, uploaded.id);
      setDocumentId(updated.document_id);
      notify("Receipt attached.");
    } catch (err) {
      notify(err instanceof ApiError ? err.message : "Couldn't attach that receipt.", "error");
    } finally {
      setAttaching(false);
    }
  }

  async function handleRemoveReceipt() {
    try {
      const updated = await expensesApi.removeReceipt(expense.id);
      setDocumentId(updated.document_id);
      notify("Receipt removed.");
    } catch (err) {
      notify(err instanceof ApiError ? err.message : "Couldn't remove that receipt.", "error");
    }
  }

  return (
    <>
      <ListItem
        secondaryAction={
          <Stack direction="row" spacing={0.5}>
            {!isViewer && <input ref={fileInputRef} type="file" hidden onChange={handleReceiptChosen} accept="image/*,.pdf" />}
            {documentId ? (
              <>
                <IconButton
                  edge="end"
                  aria-label="Download receipt"
                  onClick={() => receipt && documentsApi.download(receipt)}
                  disabled={!receipt}
                >
                  <AttachFileIcon fontSize="small" color="primary" />
                </IconButton>
                {!isViewer && (
                  <IconButton edge="end" aria-label="Remove receipt" onClick={handleRemoveReceipt}>
                    <DeleteOutlineIcon fontSize="small" />
                  </IconButton>
                )}
              </>
            ) : (
              !isViewer && (
                <IconButton
                  edge="end"
                  aria-label="Attach receipt"
                  onClick={() => fileInputRef.current?.click()}
                  disabled={attaching}
                >
                  <AttachFileIcon fontSize="small" />
                </IconButton>
              )
            )}
            <IconButton edge="end" aria-label="Show splits" onClick={toggleExpand}>
              <ExpandMoreIcon
                fontSize="small"
                sx={{ transform: expanded ? "rotate(180deg)" : "none", transition: "transform 0.15s" }}
              />
            </IconButton>
            {!isViewer && (
              <IconButton edge="end" aria-label="Delete expense" onClick={onDelete}>
                <DeleteOutlineIcon fontSize="small" />
              </IconButton>
            )}
          </Stack>
        }
      >
        <ListItemText
          primary={
            <Stack direction="row" spacing={1} sx={{ alignItems: "center", flexWrap: "wrap" }}>
              <span>{expense.description}</span>
              <Chip label={formatCents(expense.amount_cents, expense.currency)} size="small" color="primary" />
              {expense.category && <Chip label={expense.category} size="small" variant="outlined" />}
              {parentLabel && <Chip label={parentLabel} size="small" />}
              {documentId && <Chip label={receipt?.file_name ?? "Receipt attached"} size="small" variant="outlined" />}
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
              ) : isViewer ? (
                <Chip label="Unsettled" size="small" variant="outlined" />
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
