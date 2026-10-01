import { AppBar, Box, Button, Container, Stack, Toolbar, Typography } from "@mui/material";
import { Link as RouterLink, Outlet } from "react-router-dom";

const navItems = [
  { label: "Dashboard", to: "/" },
  { label: "Tasks", to: "/tasks" },
  { label: "Expenses", to: "/expenses" },
  { label: "Documents", to: "/documents" },
];

export function Layout() {
  return (
    <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      <AppBar position="static" color="primary" enableColorOnDark>
        <Toolbar>
          <Typography variant="h6" component={RouterLink} to="/" sx={{ flexGrow: 1, color: "inherit", textDecoration: "none" }}>
            Sibling Core
          </Typography>
          <Stack direction="row" spacing={1}>
            {navItems.map((item) => (
              <Button key={item.to} color="inherit" component={RouterLink} to={item.to}>
                {item.label}
              </Button>
            ))}
          </Stack>
        </Toolbar>
      </AppBar>
      <Container component="main" sx={{ flexGrow: 1, py: 4 }}>
        <Outlet />
      </Container>
    </Box>
  );
}
