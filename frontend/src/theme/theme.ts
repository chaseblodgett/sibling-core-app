import { createTheme } from "@mui/material/styles";

// MUI theme built from DESIGN.md's design tokens (Airbnb-style design
// language: white canvas, Rausch #ff385c primary, soft rounded corners,
// one shadow tier). Follow DESIGN.md for any new component or page —
// extend this theme file rather than hardcoding colors/spacing ad hoc.

const fontFamily =
  "'Airbnb Cereal VF', Circular, Inter, -apple-system, system-ui, Roboto, 'Helvetica Neue', sans-serif";

// DESIGN.md `colors`
const colors = {
  primary: "#ff385c",
  primaryActive: "#e00b41",
  primaryDisabled: "#ffd1da",
  errorText: "#c13515",
  errorTextHover: "#b32505",
  ink: "#222222",
  body: "#3f3f3f",
  muted: "#6a6a6a",
  mutedSoft: "#929292",
  hairline: "#dddddd",
  hairlineSoft: "#ebebeb",
  borderStrong: "#c1c1c1",
  canvas: "#ffffff",
  surfaceSoft: "#f7f7f7",
  surfaceStrong: "#f2f2f2",
  onPrimary: "#ffffff",
};

// DESIGN.md `rounded`
const radius = {
  xs: 4,
  sm: 8,
  md: 14,
  lg: 20,
  xl: 32,
  full: 9999,
};

// DESIGN.md Elevation — the system's single shadow tier, used sparingly
// (hover-floated cards, dropdowns). Flat (no shadow) is the baseline.
export const singleShadowTier =
  "rgba(0, 0, 0, 0.02) 0 0 0 1px, rgba(0, 0, 0, 0.04) 0 2px 6px 0, rgba(0, 0, 0, 0.1) 0 4px 8px 0";

export const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: colors.primary,
      dark: colors.primaryActive,
      light: colors.primaryDisabled,
      contrastText: colors.onPrimary,
    },
    error: {
      main: colors.errorText,
      dark: colors.errorTextHover,
    },
    text: {
      primary: colors.ink,
      secondary: colors.muted,
      disabled: colors.mutedSoft,
    },
    background: {
      default: colors.canvas,
      paper: colors.canvas,
    },
    divider: colors.hairline,
  },
  shape: {
    borderRadius: radius.sm,
  },
  typography: {
    fontFamily,
    // DESIGN.md `typography` tokens, mapped onto the closest MUI variant.
    h1: { fontSize: 28, fontWeight: 700, lineHeight: 1.43, letterSpacing: 0 }, // display-xl
    h2: { fontSize: 22, fontWeight: 500, lineHeight: 1.18, letterSpacing: "-0.44px" }, // display-lg
    h3: { fontSize: 21, fontWeight: 700, lineHeight: 1.43, letterSpacing: 0 }, // display-md
    h4: { fontSize: 20, fontWeight: 600, lineHeight: 1.2, letterSpacing: "-0.18px" }, // display-sm
    h5: { fontSize: 16, fontWeight: 600, lineHeight: 1.25, letterSpacing: 0 }, // title-md
    h6: { fontSize: 16, fontWeight: 500, lineHeight: 1.25, letterSpacing: 0 }, // title-sm
    body1: { fontSize: 16, fontWeight: 400, lineHeight: 1.5, letterSpacing: 0 }, // body-md
    body2: { fontSize: 14, fontWeight: 400, lineHeight: 1.43, letterSpacing: 0 }, // body-sm
    caption: { fontSize: 14, fontWeight: 500, lineHeight: 1.29, letterSpacing: 0 }, // caption
    button: {
      fontSize: 16,
      fontWeight: 500,
      lineHeight: 1.25,
      letterSpacing: 0,
      textTransform: "none", // DESIGN.md CTAs are normal case ("Reserve", "Continue")
    },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: { backgroundColor: colors.canvas, color: colors.ink },
      },
    },
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: ({ ownerState }) => ({
          borderRadius: radius.sm,
          height: 48,
          padding: "14px 24px",
          boxShadow: "none",
          ...(ownerState.variant === "contained" &&
            ownerState.color === "primary" && {
              backgroundColor: colors.primary,
              color: colors.onPrimary,
              "&:hover": { backgroundColor: colors.primaryActive, boxShadow: "none" },
              "&:active": { backgroundColor: colors.primaryActive },
              "&.Mui-disabled": { backgroundColor: colors.primaryDisabled, color: colors.onPrimary },
            }),
        }),
        outlined: {
          borderColor: colors.ink,
          color: colors.ink,
          "&:hover": { borderColor: colors.ink, backgroundColor: colors.surfaceSoft },
        },
        text: {
          color: colors.ink,
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: radius.sm,
          backgroundColor: colors.canvas,
          "& .MuiOutlinedInput-notchedOutline": { borderColor: colors.hairline },
          "&:hover .MuiOutlinedInput-notchedOutline": { borderColor: colors.borderStrong },
          "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
            borderColor: colors.ink,
            borderWidth: 2,
          },
          "&.Mui-error .MuiOutlinedInput-notchedOutline": { borderColor: colors.errorText },
        },
        input: { padding: "14px 12px" },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: { color: colors.muted, "&.Mui-focused": { color: colors.ink } },
      },
    },
    MuiFormHelperText: {
      styleOverrides: {
        root: { marginLeft: 2 },
      },
    },
    MuiPaper: {
      defaultProps: { elevation: 0 },
      styleOverrides: {
        root: { backgroundImage: "none" },
        rounded: { borderRadius: radius.md },
      },
    },
    MuiAlert: {
      styleOverrides: {
        root: { borderRadius: radius.sm },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: colors.canvas,
          color: colors.ink,
          borderBottom: `1px solid ${colors.hairline}`,
        },
      },
    },
    MuiLink: {
      styleOverrides: {
        root: { color: colors.ink, fontWeight: 500 },
      },
    },
  },
});
