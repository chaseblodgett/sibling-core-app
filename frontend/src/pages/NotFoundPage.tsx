import { Box, Button, Typography } from "@mui/material";
import { Link as RouterLink } from "react-router-dom";

export function NotFoundPage() {
  return (
    <Box sx={{ textAlign: "center", mt: { xs: 6, sm: 10 } }}>
      <Typography variant="h2" sx={{ fontWeight: 600, mb: 1 }}>
        404
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
        That page doesn't exist.
      </Typography>
      <Button variant="contained" component={RouterLink} to="/">
        Back to Dashboard
      </Button>
    </Box>
  );
}
