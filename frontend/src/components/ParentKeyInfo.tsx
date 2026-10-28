import { Box, Button, Stack, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { ApiError } from "../api/client";
import { parentsApi } from "../api/parents";
import { useNotify } from "../context/NotificationContext";
import type { Parent, ParentKeyInfo as ParentKeyInfoFields } from "../types";

interface ParentKeyInfoProps {
  parent: Parent;
  onSaved: (updated: Parent) => void;
  readOnly?: boolean;
}

function fieldsFrom(parent: Parent): ParentKeyInfoFields {
  return {
    doctor_name: parent.doctor_name,
    doctor_phone: parent.doctor_phone,
    allergies: parent.allergies,
    medications: parent.medications,
    insurance_provider: parent.insurance_provider,
    insurance_member_id: parent.insurance_member_id,
    emergency_contact_name: parent.emergency_contact_name,
    emergency_contact_phone: parent.emergency_contact_phone,
  };
}

// Freeform text, not a structured medical record — the DESIGN.md
// card-tint-yellow-bold token, reserved for high-emphasis content, gives
// this its own visual identity so it reads as "important" rather than
// blending into the rest of the parent's row.
const KEY_INFO_BG = "#f9e79f";
const KEY_INFO_BORDER = "#dd5b00";

export function ParentKeyInfo({ parent, onSaved, readOnly = false }: ParentKeyInfoProps) {
  const notify = useNotify();
  const [form, setForm] = useState<ParentKeyInfoFields>(fieldsFrom(parent));
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  function set<K extends keyof ParentKeyInfoFields>(key: K, value: string) {
    setForm((prev) => ({ ...prev, [key]: value }));
  }

  async function handleSave() {
    setError(null);
    setSaving(true);
    try {
      const updated = await parentsApi.update(parent.id, form);
      onSaved(updated);
      notify("Key info updated.");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Couldn't save key info.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <Box sx={{ bgcolor: KEY_INFO_BG, border: `1px solid ${KEY_INFO_BORDER}`, borderRadius: 2, p: 2.5, mt: 1.5 }}>
      <Typography variant="subtitle2" gutterBottom>
        Key info for {parent.name}
      </Typography>
      <Typography variant="caption" color="text.secondary" sx={{ display: "block", mb: 2 }}>
        For a sibling who needs this fast — stepping in during an emergency, or on a call with a new doctor.
      </Typography>
      {error && (
        <Typography variant="body2" color="error" sx={{ mb: 1.5 }}>
          {error}
        </Typography>
      )}
      <Stack spacing={2}>
        <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap" }}>
          <TextField
            label="Doctor name"
            size="small"
            value={form.doctor_name}
            onChange={(e) => set("doctor_name", e.target.value)}
            sx={{ flex: 1, minWidth: 180, bgcolor: "background.paper" }}
            slotProps={{ input: { readOnly } }}
          />
          <TextField
            label="Doctor phone"
            size="small"
            value={form.doctor_phone}
            onChange={(e) => set("doctor_phone", e.target.value)}
            sx={{ flex: 1, minWidth: 180, bgcolor: "background.paper" }}
            slotProps={{ input: { readOnly } }}
          />
        </Stack>
        <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap" }}>
          <TextField
            label="Emergency contact name"
            size="small"
            value={form.emergency_contact_name}
            onChange={(e) => set("emergency_contact_name", e.target.value)}
            sx={{ flex: 1, minWidth: 180, bgcolor: "background.paper" }}
            slotProps={{ input: { readOnly } }}
          />
          <TextField
            label="Emergency contact phone"
            size="small"
            value={form.emergency_contact_phone}
            onChange={(e) => set("emergency_contact_phone", e.target.value)}
            sx={{ flex: 1, minWidth: 180, bgcolor: "background.paper" }}
            slotProps={{ input: { readOnly } }}
          />
        </Stack>
        <Stack direction="row" spacing={2} sx={{ flexWrap: "wrap" }}>
          <TextField
            label="Insurance provider"
            size="small"
            value={form.insurance_provider}
            onChange={(e) => set("insurance_provider", e.target.value)}
            sx={{ flex: 1, minWidth: 180, bgcolor: "background.paper" }}
            slotProps={{ input: { readOnly } }}
          />
          <TextField
            label="Insurance member ID"
            size="small"
            value={form.insurance_member_id}
            onChange={(e) => set("insurance_member_id", e.target.value)}
            sx={{ flex: 1, minWidth: 180, bgcolor: "background.paper" }}
            slotProps={{ input: { readOnly } }}
          />
        </Stack>
        <TextField
          label="Allergies"
          size="small"
          value={form.allergies}
          onChange={(e) => set("allergies", e.target.value)}
          multiline
          minRows={2}
          fullWidth
          sx={{ bgcolor: "background.paper" }}
          slotProps={{ input: { readOnly } }}
        />
        <TextField
          label="Medications"
          size="small"
          value={form.medications}
          onChange={(e) => set("medications", e.target.value)}
          multiline
          minRows={2}
          fullWidth
          helperText="A readable list for someone else, not a dosage chart."
          sx={{ bgcolor: "background.paper" }}
          slotProps={{ input: { readOnly } }}
        />
        {!readOnly && (
          <Box>
            <Button variant="contained" size="small" onClick={handleSave} disabled={saving}>
              Save key info
            </Button>
          </Box>
        )}
      </Stack>
    </Box>
  );
}
