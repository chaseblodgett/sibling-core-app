import { Alert, Box, Button, Divider, Link, Paper, Stack, TextField, Typography } from "@mui/material";
import { type FormEvent, useState } from "react";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import { ApiError } from "../api/client";
import { GoogleSignInButton } from "../components/GoogleSignInButton";
import { useAuth } from "../context/AuthContext";
import { cardRadius, singleShadowTier } from "../theme/theme";

const EMAIL_PATTERN = /^\S+@\S+\.\S+$/;
const MIN_PASSWORD_LENGTH = 8;

interface FieldErrors {
  name?: string;
  email?: string;
  password?: string;
}

export function SignupPage() {
  const { signUp, loginWithGoogle } = useAuth();
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleGoogleIdToken(idToken: string) {
    setServerError(null);
    try {
      await loginWithGoogle(idToken);
      navigate("/app", { replace: true });
    } catch (err) {
      setServerError(err instanceof ApiError ? err.message : "Google sign-in failed. Try again.");
    }
  }

  function validate(): boolean {
    const errors: FieldErrors = {};
    if (name.trim().length === 0) {
      errors.name = "Name is required";
    }
    if (!EMAIL_PATTERN.test(email)) {
      errors.email = "Enter a valid email address";
    }
    if (password.length < MIN_PASSWORD_LENGTH) {
      errors.password = `Password must be at least ${MIN_PASSWORD_LENGTH} characters`;
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
      await signUp(email, password, name);
      navigate("/app", { replace: true });
    } catch (err) {
      setServerError(err instanceof ApiError ? err.message : "Something went wrong. Try again.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Box sx={{ display: "flex", justifyContent: "center", mt: { xs: 4, sm: 8 }, px: 2 }}>
      <Paper
        sx={{ p: 4, width: "100%", maxWidth: 400, borderRadius: `${cardRadius}px`, boxShadow: singleShadowTier }}
      >
        <Stack spacing={3} component="form" onSubmit={handleSubmit} noValidate>
          <Box>
            <Typography variant="h4" component="h1" gutterBottom>
              Sign up
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Create an account to get your family organized.
            </Typography>
          </Box>

          {serverError && <Alert severity="error">{serverError}</Alert>}

          <GoogleSignInButton text="signup_with" onIdToken={handleGoogleIdToken} />

          <Divider>or</Divider>

          <TextField
            label="Full name"
            name="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            error={!!fieldErrors.name}
            helperText={fieldErrors.name}
            autoComplete="name"
            fullWidth
          />
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
            helperText={fieldErrors.password ?? `At least ${MIN_PASSWORD_LENGTH} characters`}
            autoComplete="new-password"
            fullWidth
          />
          <Button type="submit" variant="contained" size="large" disabled={submitting} fullWidth>
            {submitting ? "Creating account…" : "Create account"}
          </Button>
          <Typography variant="body2">
            Already have an account? <Link component={RouterLink} to="/login">Log in</Link>
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
}
