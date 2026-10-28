import { Box, CircularProgress, FormControlLabel, Paper, Switch, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { notificationPreferencesApi } from "../api/notificationPreferences";
import { ApiError } from "../api/client";
import { useNotify } from "../context/NotificationContext";

export function SettingsPage() {
  const notify = useNotify();
  const [taskAssignedEmail, setTaskAssignedEmail] = useState(true);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    notificationPreferencesApi
      .get()
      .then((prefs) => setTaskAssignedEmail(prefs.task_assigned_email))
      .catch(() => {
        // default (on) stays in place — the settings toggle still works,
        // it just can't show what's actually saved
      })
      .finally(() => setLoading(false));
  }, []);

  async function handleToggle(checked: boolean) {
    setTaskAssignedEmail(checked);
    setSaving(true);
    try {
      await notificationPreferencesApi.update(checked);
      notify(checked ? "You'll get an email when you're assigned a task." : "Task-assignment emails turned off.");
    } catch (err) {
      setTaskAssignedEmail(!checked);
      notify(err instanceof ApiError ? err.message : "Couldn't save that setting.", "error");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", mt: 8 }}>
        <CircularProgress sx={{ color: "primary.main" }} />
      </Box>
    );
  }

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Settings
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Notification preferences apply to every family you're part of.
      </Typography>

      <Paper sx={{ p: 3, maxWidth: 480 }}>
        <Typography variant="h6" component="h2" gutterBottom>
          Email notifications
        </Typography>
        <FormControlLabel
          control={
            <Switch
              checked={taskAssignedEmail}
              onChange={(e) => handleToggle(e.target.checked)}
              disabled={saving}
            />
          }
          label="Email me when I'm assigned a task"
        />
      </Paper>
    </>
  );
}
