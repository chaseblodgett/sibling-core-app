import AddIcon from "@mui/icons-material/Add";
import { Box, Button, Grid, List, ListItem, ListItemText, Paper, Typography } from "@mui/material";

export function ExpensesPage() {
  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
        <Typography variant="h4" component="h1">
          Expenses
        </Typography>
        <Button variant="contained" startIcon={<AddIcon />}>
          Log expense
        </Button>
      </Box>
      <Grid container spacing={2}>
        <Grid size={{ xs: 12, md: 8 }}>
          <Paper>
            <List>
              <ListItem>
                <ListItemText primary="No expenses yet" secondary="Logged expenses and their splits will show up here." />
              </ListItem>
            </List>
          </Paper>
        </Grid>
        <Grid size={{ xs: 12, md: 4 }}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Balances
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Running balances per member will appear here once expenses are logged.
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </>
  );
}
