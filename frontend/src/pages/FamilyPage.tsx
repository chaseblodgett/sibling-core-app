import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  Grid,
  IconButton,
  InputAdornment,
  List,
  ListItem,
  ListItemText,
  Paper,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { type FormEvent, useCallback, useEffect, useState } from "react";
import { ApiError } from "../api/client";
import { familyApi } from "../api/family";
import { useAuth } from "../context/AuthContext";
import { useFamilies } from "../context/FamilyContext";
import type { Family, FamilyInvite, Member } from "../types";

export function FamilyPage() {
  const { families, loading: familiesLoading, reload } = useFamilies();

  const [createName, setCreateName] = useState("");
  const [joinCode, setJoinCode] = useState("");
  const [onboardSubmitting, setOnboardSubmitting] = useState(false);
  const [onboardError, setOnboardError] = useState<string | null>(null);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setOnboardError(null);
    setOnboardSubmitting(true);
    try {
      await familyApi.create(createName);
      setCreateName("");
      await reload();
    } catch (err) {
      setOnboardError(err instanceof ApiError ? err.message : "Couldn't create the family.");
    } finally {
      setOnboardSubmitting(false);
    }
  }

  async function handleJoin(e: FormEvent) {
    e.preventDefault();
    setOnboardError(null);
    setOnboardSubmitting(true);
    try {
      await familyApi.join(joinCode);
      setJoinCode("");
      await reload();
    } catch (err) {
      setOnboardError(err instanceof ApiError ? err.message : "Couldn't join with that invite code.");
    } finally {
      setOnboardSubmitting(false);
    }
  }

  if (familiesLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 8 }}>
        <CircularProgress sx={{ color: "primary.main" }} />
      </Box>
    );
  }

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Family
      </Typography>

      {onboardError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {onboardError}
        </Alert>
      )}

      <Grid container spacing={2} sx={{ mb: 2 }}>
        <Grid size={{ xs: 12, sm: 6 }}>
          <Paper sx={{ p: 3 }}>
            <Stack spacing={2} component="form" onSubmit={handleCreate}>
              <Box>
                <Typography variant="h6" component="h2">
                  Create a family
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Start a new family and invite siblings to join.
                </Typography>
              </Box>
              <TextField
                label="Family name"
                value={createName}
                onChange={(e) => setCreateName(e.target.value)}
                required
                fullWidth
              />
              <Button type="submit" variant="contained" disabled={onboardSubmitting}>
                Create family
              </Button>
            </Stack>
          </Paper>
        </Grid>
        <Grid size={{ xs: 12, sm: 6 }}>
          <Paper sx={{ p: 3 }}>
            <Stack spacing={2} component="form" onSubmit={handleJoin}>
              <Box>
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
              <Button type="submit" variant="outlined" disabled={onboardSubmitting}>
                Join family
              </Button>
            </Stack>
          </Paper>
        </Grid>
      </Grid>

      {families.length === 0 ? (
        <Alert severity="info">You're not part of any family yet — create or join one above.</Alert>
      ) : (
        <Stack spacing={3}>
          {families.map((family) => (
            <FamilyCard key={family.id} family={family} />
          ))}
        </Stack>
      )}
    </>
  );
}

function FamilyCard({ family }: { family: Family }) {
  const { user } = useAuth();

  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [inviteEmail, setInviteEmail] = useState("");
  const [inviteSubmitting, setInviteSubmitting] = useState(false);
  const [sentInvites, setSentInvites] = useState<FamilyInvite[]>([]);
  const [copied, setCopied] = useState(false);

  const loadMembers = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      setMembers(await familyApi.listMembers(family.id));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't load members.");
    } finally {
      setLoading(false);
    }
  }, [family.id]);

  useEffect(() => {
    loadMembers();
  }, [loadMembers]);

  async function handleInvite(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setInviteSubmitting(true);
    try {
      const invite = await familyApi.invite(family.id, inviteEmail);
      setSentInvites((prev) => [invite, ...prev]);
      setInviteEmail("");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't send that invite.");
    } finally {
      setInviteSubmitting(false);
    }
  }

  async function handleRemove(memberId: string) {
    setError(null);
    try {
      await familyApi.removeMember(family.id, memberId);
      setMembers((prev) => prev.filter((m) => m.id !== memberId));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't remove that member.");
    }
  }

  function handleCopyInviteCode() {
    navigator.clipboard.writeText(family.invite_code).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  }

  const isOwner = members.find((m) => m.user_id === user?.id)?.role === "owner";

  return (
    <Stack spacing={2}>
      {error && <Alert severity="error">{error}</Alert>}

      <Paper sx={{ p: 3 }}>
        <Stack direction="row" sx={{ justifyContent: "space-between", alignItems: "flex-start", mb: 2 }}>
          <Box>
            <Typography variant="h6" component="h2">
              {family.name}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {loading ? "Loading members…" : `${members.length} ${members.length === 1 ? "member" : "members"}`}
            </Typography>
          </Box>
          {isOwner && <Chip label="Owner" color="primary" size="small" />}
        </Stack>
        <TextField
          label="Invite code"
          value={family.invite_code}
          fullWidth
          slotProps={{
            input: {
              readOnly: true,
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton onClick={handleCopyInviteCode} edge="end" aria-label="Copy invite code">
                    <ContentCopyIcon fontSize="small" />
                  </IconButton>
                </InputAdornment>
              ),
            },
          }}
          helperText={copied ? "Copied to clipboard" : "Share this code so a sibling can join"}
        />
      </Paper>

      {isOwner && (
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" component="h2" gutterBottom>
            Invite by email
          </Typography>
          <Stack direction="row" spacing={2} component="form" onSubmit={handleInvite}>
            <TextField
              label="Email"
              type="email"
              value={inviteEmail}
              onChange={(e) => setInviteEmail(e.target.value)}
              required
              fullWidth
              size="small"
            />
            <Button type="submit" variant="contained" disabled={inviteSubmitting}>
              Send invite
            </Button>
          </Stack>
          {sentInvites.length > 0 && (
            <List dense sx={{ mt: 1 }}>
              {sentInvites.map((invite) => (
                <ListItem key={invite.id} disableGutters>
                  <ListItemText primary={invite.email} secondary="Invited — share the code above with them too" />
                </ListItem>
              ))}
            </List>
          )}
        </Paper>
      )}

      <Paper>
        <List>
          {members.map((member) => (
            <ListItem
              key={member.id}
              secondaryAction={
                isOwner &&
                member.user_id !== user?.id && (
                  <IconButton edge="end" aria-label="Remove member" onClick={() => handleRemove(member.id)}>
                    <DeleteOutlineIcon fontSize="small" />
                  </IconButton>
                )
              }
            >
              <ListItemText
                primary={member.full_name || member.email}
                secondary={member.role === "owner" ? "Owner" : "Member"}
              />
            </ListItem>
          ))}
        </List>
      </Paper>
    </Stack>
  );
}
