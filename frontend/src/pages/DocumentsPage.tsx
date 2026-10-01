import UploadIcon from "@mui/icons-material/Upload";
import { Box, Button, List, ListItem, ListItemText, Paper, Typography } from "@mui/material";

export function DocumentsPage() {
  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
        <Typography variant="h4" component="h1">
          Documents
        </Typography>
        <Button variant="contained" startIcon={<UploadIcon />}>
          Upload
        </Button>
      </Box>
      <Paper>
        <List>
          <ListItem>
            <ListItemText primary="No documents yet" secondary="Files uploaded for your family will show up here." />
          </ListItem>
        </List>
      </Paper>
    </>
  );
}
