import { AppBar, Box, Button, Container, MenuItem, Select, Stack, Toolbar, Typography } from "@mui/material";
import { Link as RouterLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useFamilies } from "../context/FamilyContext";

const navItems = [
  { label: "Dashboard", to: "/" },
  { label: "Family", to: "/family" },
  { label: "Parents", to: "/parents" },
  { label: "Tasks", to: "/tasks" },
  { label: "Expenses", to: "/expenses" },
  { label: "Documents", to: "/documents" },
];

export function Layout() {
  const { logout } = useAuth();
  const { families, selectedFamilyId, setSelectedFamilyId } = useFamilies();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login", { replace: true });
  }

  return (
    <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      <AppBar position="static" color="transparent" elevation={0}>
        <Toolbar>
          <Typography
            variant="h6"
            component={RouterLink}
            to="/"
            sx={{ flexGrow: 1, color: "inherit", textDecoration: "none", fontWeight: 700 }}
          >
            Caretree
          </Typography>
          <Stack direction="row" spacing={1} sx={{ alignItems: "center" }}>
            {navItems.map((item) => (
              <Button key={item.to} color="inherit" component={RouterLink} to={item.to}>
                {item.label}
              </Button>
            ))}
            <Button variant="outlined" onClick={handleLogout}>
              Log out
            </Button>
          </Stack>
        </Toolbar>
        {families.length > 0 && (
          <Toolbar variant="dense" sx={{ borderTop: "1px solid", borderColor: "divider", minHeight: 44 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mr: 1.5 }}>
              Viewing:
            </Typography>
            {families.length === 1 ? (
              <Typography variant="body2" sx={{ fontWeight: 500 }}>
                {families[0].name}
              </Typography>
            ) : (
              <Select
                size="small"
                value={selectedFamilyId ?? ""}
                onChange={(e) => setSelectedFamilyId(e.target.value || null)}
                sx={{ minWidth: 200 }}
              >
                <MenuItem value="">All families</MenuItem>
                {families.map((f) => (
                  <MenuItem key={f.id} value={f.id}>
                    {f.name}
                  </MenuItem>
                ))}
              </Select>
            )}
          </Toolbar>
        )}
      </AppBar>
      <Container component="main" sx={{ flexGrow: 1, py: 4 }}>
        <Outlet />
      </Container>
    </Box>
  );
}
