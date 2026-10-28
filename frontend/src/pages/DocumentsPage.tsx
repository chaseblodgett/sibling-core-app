import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import DownloadIcon from "@mui/icons-material/Download";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import UploadIcon from "@mui/icons-material/Upload";
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
  Select,
  Stack,
  Typography,
} from "@mui/material";
import { type ChangeEvent, useCallback, useEffect, useRef, useState } from "react";
import { ApiError } from "../api/client";
import { documentsApi } from "../api/documents";
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { EmptyState } from "../components/EmptyState";
import { FamilySection } from "../components/FamilySection";
import { useFamilies } from "../context/FamilyContext";
import { useNotify } from "../context/NotificationContext";
import { useConfirm } from "../hooks/useConfirm";
import { familyColorMap } from "../utils/familyColor";
import type { Document, Family, Member, Parent } from "../types";

function formatBytes(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function DocumentsPage() {
  const { families, loading: familiesLoading } = useFamilies();
  const notify = useNotify();
  const { state: confirmState, confirm, close: closeConfirm, handleConfirm } = useConfirm();

  const [documents, setDocuments] = useState<Document[]>([]);
  const [parents, setParents] = useState<Parent[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [uploadFamilyId, setUploadFamilyId] = useState("");
  const [uploadParentId, setUploadParentId] = useState("");
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const load = useCallback(async (targets: Family[]) => {
    if (targets.length === 0) {
      setDocuments([]);
      setParents([]);
      setMembers([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const [documentLists, parentLists, memberLists] = await Promise.all([
        Promise.all(targets.map((f) => documentsApi.list(f.id))),
        Promise.all(targets.map((f) => parentsApi.list(f.id))),
        Promise.all(targets.map((f) => familyApi.listMembers(f.id))),
      ]);
      setDocuments(documentLists.flat());
      setParents(parentLists.flat());
      setMembers(memberLists.flat());
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't load documents.");
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
    setUploadFamilyId((prev) => prev || families[0]?.id || "");
  }, [families]);

  const uploadFamilyParents = parents.filter((p) => p.family_id === uploadFamilyId);

  async function handleFileChosen(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    e.target.value = "";
    if (!file || !uploadFamilyId) return;

    setError(null);
    setUploading(true);
    try {
      await documentsApi.upload(uploadFamilyId, file, uploadParentId || undefined);
      setUploadParentId("");
      notify(`${file.name} uploaded.`);
      await load(families);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't upload that file.");
    } finally {
      setUploading(false);
    }
  }

  async function handleDownload(doc: Document) {
    setError(null);
    try {
      await documentsApi.download(doc);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't download that file.");
    }
  }

  function confirmDelete(doc: Document) {
    confirm({
      title: "Delete document?",
      message: `Delete "${doc.file_name}"? This can't be undone.`,
      confirmLabel: "Delete",
      destructive: true,
      onConfirm: () => handleDelete(doc.id),
    });
  }

  async function handleDelete(docId: string) {
    setError(null);
    try {
      await documentsApi.delete(docId);
      setDocuments((prev) => prev.filter((d) => d.id !== docId));
      notify("Document deleted.");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't delete that document.");
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
        icon={<InsertDriveFileIcon />}
        title="No family yet"
        description="Create or join a family first, then come back here to start sharing documents."
      />
    );
  }

  const colors = familyColorMap(families);

  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2, flexWrap: "wrap", gap: 2 }}>
        <Typography variant="h4" component="h1">
          Documents
        </Typography>
        <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
          {families.length > 1 && (
            <Select size="small" value={uploadFamilyId} onChange={(e) => setUploadFamilyId(e.target.value)}>
              {families.map((f) => (
                <MenuItem key={f.id} value={f.id}>
                  {f.name}
                </MenuItem>
              ))}
            </Select>
          )}
          <Select
            size="small"
            value={uploadParentId}
            onChange={(e) => setUploadParentId(e.target.value)}
            displayEmpty
            sx={{ minWidth: 160 }}
          >
            <MenuItem value="">No parent (general)</MenuItem>
            {uploadFamilyParents.map((p) => (
              <MenuItem key={p.id} value={p.id}>
                {p.name}
              </MenuItem>
            ))}
          </Select>
          <input ref={fileInputRef} type="file" hidden onChange={handleFileChosen} />
          <Button
            variant="contained"
            startIcon={<UploadIcon />}
            disabled={uploading || !uploadFamilyId}
            onClick={() => fileInputRef.current?.click()}
          >
            {uploading ? "Uploading…" : "Upload"}
          </Button>
        </Stack>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Stack spacing={2}>
        {families.map((family) => {
          const familyDocuments = documents.filter((d) => d.family_id === family.id);
          return (
            <FamilySection key={family.id} family={family} color={colors.get(family.id)!}>
              {familyDocuments.length === 0 ? (
                <EmptyState
                  icon={<InsertDriveFileIcon />}
                  title="No documents yet"
                  description="Files uploaded for this family will show up here — upload one using the button above."
                />
              ) : (
                <List>
                  {familyDocuments.map((doc) => (
                    <ListItem
                      key={doc.id}
                      secondaryAction={
                        <Stack direction="row" spacing={0.5}>
                          <IconButton edge="end" aria-label="Download" onClick={() => handleDownload(doc)}>
                            <DownloadIcon fontSize="small" />
                          </IconButton>
                          <IconButton edge="end" aria-label="Delete document" onClick={() => confirmDelete(doc)}>
                            <DeleteOutlineIcon fontSize="small" />
                          </IconButton>
                        </Stack>
                      }
                    >
                      <ListItemText
                        primary={
                          <Stack direction="row" spacing={1} sx={{ alignItems: "center", flexWrap: "wrap" }}>
                            <span>{doc.file_name}</span>
                            <Chip label={formatBytes(doc.size_bytes)} size="small" variant="outlined" />
                            {parentName(doc.parent_id) && <Chip label={parentName(doc.parent_id)} size="small" />}
                          </Stack>
                        }
                        secondary={`Uploaded by ${memberName(doc.uploaded_by) ?? "someone"}`}
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
