import { createTheme } from "@mui/material/styles";

// Basic MUI theme. Adjust palette/typography as the product direction firms up.
export const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#3f51b5",
    },
    secondary: {
      main: "#00897b",
    },
  },
  shape: {
    borderRadius: 8,
  },
});
