import AddIcon from "@mui/icons-material/Add";
import { Box, Button, List, ListItem, ListItemText, Paper, Typography } from "@mui/material";

export function TasksPage() {
  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
        <Typography variant="h4" component="h1">
          Tasks
        </Typography>
        <Button variant="contained" startIcon={<AddIcon />}>
          New task
        </Button>
      </Box>
      <Paper>
        <List>
          <ListItem>
            <ListItemText primary="No tasks yet" secondary="Tasks assigned to family members will show up here." />
          </ListItem>
        </List>
      </Paper>
    </>
  );
}
