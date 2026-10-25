import { Card, CardContent, Grid, Typography } from "@mui/material";
import { useAuth } from "../context/AuthContext";

export function DashboardPage() {
  const { user } = useAuth();

  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Welcome{user ? `, ${user.name}` : ""}
      </Typography>
      <Grid container spacing={2}>
        <Grid size={{ xs: 12, sm: 4 }}>
          <Card>
            <CardContent>
              <Typography variant="overline">Open tasks</Typography>
              <Typography variant="h4">—</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid size={{ xs: 12, sm: 4 }}>
          <Card>
            <CardContent>
              <Typography variant="overline">Outstanding balance</Typography>
              <Typography variant="h4">—</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid size={{ xs: 12, sm: 4 }}>
          <Card>
            <CardContent>
              <Typography variant="overline">Documents</Typography>
              <Typography variant="h4">—</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </>
  );
}
