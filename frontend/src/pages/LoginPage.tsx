import { Alert, Box, Button, Link, Paper, Stack, TextField, Typography } from "@mui/material";
import { type FormEvent, useState } from "react";
import { Link as RouterLink, useLocation, useNavigate } from "react-router-dom";
import type { Location } from "react-router-dom";
import { ApiError } from "../api/client";
import { useAuth } from "../context/AuthContext";
import { singleShadowTier } from "../theme/theme";

const EMAIL_PATTERN = /^\S+@\S+\.\S+$/;

interface FieldErrors {
  email?: string;
  password?: string;
}

export function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  function validate(): boolean {
    const errors: FieldErrors = {};
    if (!EMAIL_PATTERN.test(email)) {
      errors.email = "Enter a valid email address";
    }
    if (password.length === 0) {
      errors.password = "Password is required";
    }
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setServerError(null);
    if (!validate()) {
      return;
    }

    setSubmitting(true);
    try {
      await login(email, password);
      const from = (location.state as { from?: Pick<Location, "pathname"> } | null)?.from;
      navigate(from?.pathname ?? "/", { replace: true });
    } catch (err) {
      setServerError(err instanceof ApiError ? err.message : "Something went wrong. Try again.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Box sx={{ display: "flex", justifyContent: "center", mt: { xs: 4, sm: 8 }, px: 2 }}>
      <Paper
        sx={{ p: 4, width: "100%", maxWidth: 400, borderRadius: "14px", boxShadow: singleShadowTier }}
      >
        <Stack spacing={3} component="form" onSubmit={handleSubmit} noValidate>
          <Box>
            <Typography variant="h4" component="h1" gutterBottom>
              Log in
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Welcome back — sign in to your family.
            </Typography>
          </Box>

          {serverError && <Alert severity="error">{serverError}</Alert>}

          <TextField
            label="Email"
            type="email"
            name="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            error={!!fieldErrors.email}
            helperText={fieldErrors.email}
            autoComplete="email"
            fullWidth
          />
          <TextField
            label="Password"
            type="password"
            name="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            error={!!fieldErrors.password}
            helperText={fieldErrors.password}
            autoComplete="current-password"
            fullWidth
          />
          <Button type="submit" variant="contained" size="large" disabled={submitting} fullWidth>
            {submitting ? "Logging in…" : "Log in"}
          </Button>
          <Typography variant="body2">
            Don&apos;t have an account? <Link component={RouterLink} to="/signup">Sign up</Link>
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
}
