import DeleteOutlineIcon from "@mui/icons-material/DeleteOutlined";
import DownloadIcon from "@mui/icons-material/Download";
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
  Paper,
  Select,
  Stack,
  Typography,
} from "@mui/material";
import { type ChangeEvent, useCallback, useEffect, useRef, useState } from "react";
import { ApiError } from "../api/client";
import { documentsApi } from "../api/documents";
import { familyApi } from "../api/family";
import { parentsApi } from "../api/parents";
import { useFamilies } from "../context/FamilyContext";
import type { Document, Member, Parent } from "../types";

function formatBytes(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function DocumentsPage() {
  const { families, loading: familiesLoading, selectedFamilyId } = useFamilies();

  const [documents, setDocuments] = useState<Document[]>([]);
  const [parents, setParents] = useState<Parent[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [uploadFamilyId, setUploadFamilyId] = useState("");
  const [uploadParentId, setUploadParentId] = useState("");
  const [uploading, setUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const viewingAllFamilies = selectedFamilyId === null;
  const familiesToLoad = viewingAllFamilies ? families : families.filter((f) => f.id === selectedFamilyId);

  const load = useCallback(async (targets: typeof families) => {
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
      load(familiesToLoad);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [familiesLoading, selectedFamilyId, families]);

  useEffect(() => {
    setUploadFamilyId(selectedFamilyId ?? families[0]?.id ?? "");
  }, [selectedFamilyId, families]);

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
      await load(familiesToLoad);
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

  async function handleDelete(docId: string) {
    setError(null);
    try {
      await documentsApi.delete(docId);
      setDocuments((prev) => prev.filter((d) => d.id !== docId));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't delete that document.");
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
    return <Alert severity="info">You need a family first to share documents.</Alert>;
  }

  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2, flexWrap: "wrap", gap: 2 }}>
        <Typography variant="h4" component="h1">
          Documents
        </Typography>
        <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
          {viewingAllFamilies && families.length > 1 && (
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

      <Paper>
        <List>
          {documents.length === 0 && (
            <ListItem>
              <ListItemText primary="No documents yet" secondary="Files uploaded for your family will show up here." />
            </ListItem>
          )}
          {documents.map((doc) => (
            <ListItem
              key={doc.id}
              secondaryAction={
                <Stack direction="row" spacing={0.5}>
                  <IconButton edge="end" aria-label="Download" onClick={() => handleDownload(doc)}>
                    <DownloadIcon fontSize="small" />
                  </IconButton>
                  <IconButton edge="end" aria-label="Delete document" onClick={() => handleDelete(doc.id)}>
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
                    {viewingAllFamilies && (
                      <Chip label={familyName(doc.family_id)} size="small" color="primary" variant="outlined" />
                    )}
                    {parentName(doc.parent_id) && <Chip label={parentName(doc.parent_id)} size="small" />}
                  </Stack>
                }
                secondary={`Uploaded by ${memberName(doc.uploaded_by) ?? "someone"}`}
              />
            </ListItem>
          ))}
        </List>
      </Paper>
    </>
  );
}
