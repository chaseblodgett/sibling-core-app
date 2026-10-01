import { Box, Button, Link, Paper, Stack, TextField, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

export function LoginPage() {
  return (
    <Box sx={{ display: "flex", justifyContent: "center", mt: 6 }}>
      <Paper sx={{ p: 4, width: "100%", maxWidth: 400 }} elevation={2}>
        <Stack spacing={2} component="form">
          <Typography variant="h5" component="h1">
            Log in
          </Typography>
          <TextField label="Email" type="email" name="email" fullWidth />
          <TextField label="Password" type="password" name="password" fullWidth />
          <Button type="submit" variant="contained" size="large">
            Log in
          </Button>
          <Typography variant="body2">
            Don&apos;t have an account? <Link component={RouterLink} to="/signup">Sign up</Link>
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
}
