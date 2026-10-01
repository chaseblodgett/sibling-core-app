import { Box, Button, Link, Paper, Stack, TextField, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

export function SignupPage() {
  return (
    <Box sx={{ display: "flex", justifyContent: "center", mt: 6 }}>
      <Paper sx={{ p: 4, width: "100%", maxWidth: 400 }} elevation={2}>
        <Stack spacing={2} component="form">
          <Typography variant="h5" component="h1">
            Sign up
          </Typography>
          <TextField label="Full name" name="fullName" fullWidth />
          <TextField label="Email" type="email" name="email" fullWidth />
          <TextField label="Password" type="password" name="password" fullWidth />
          <Button type="submit" variant="contained" size="large">
            Create account
          </Button>
          <Typography variant="body2">
            Already have an account? <Link component={RouterLink} to="/login">Log in</Link>
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
}
