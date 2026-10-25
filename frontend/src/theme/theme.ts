import { createTheme } from "@mui/material/styles";

// MUI theme built from DESIGN.md's design tokens (Notion-style design
// language: white canvas, signature purple #5645d4 primary, rectangular
// 8px buttons/inputs, 12px-rounded cards, link-blue reserved for inline
// links). Follow DESIGN.md for any new component or page — extend this
// theme file rather than hardcoding colors/spacing ad hoc.

// Notion Sans is proprietary; DESIGN.md documents it as Inter-based, so
// Inter (already loaded in index.html) is the open-source substitute —
// same substitution pattern used for the previous Airbnb Cereal system.
const fontFamily =
  "'Inter', -apple-system, system-ui, 'Segoe UI', Helvetica, sans-serif";

// DESIGN.md `colors`
const colors = {
  primary: "#5645d4",
  primaryPressed: "#4534b3",
  onPrimary: "#ffffff",
  linkBlue: "#0075de",
  linkBluePressed: "#005bab",
  canvas: "#ffffff",
  surface: "#f6f5f4",
  surfaceSoft: "#fafaf9",
  hairline: "#e5e3df",
  hairlineSoft: "#ede9e4",
  hairlineStrong: "#c8c4be",
  inkDeep: "#000000",
  ink: "#1a1a1a",
  charcoal: "#37352f",
  slate: "#5d5b54",
  steel: "#787671",
  muted: "#bbb8b1",
  onDark: "#ffffff",
  success: "#1aae39",
  warning: "#dd5b00",
  error: "#e03131",
};

// DESIGN.md `rounded` — buttons/inputs use md (rectangular, NOT pills);
// cards use lg; full is reserved for pill tabs/badges only.
const radius = {
  xs: 4,
  sm: 6,
  md: 8,
  lg: 12,
  xl: 16,
  full: 9999,
};

// DESIGN.md Elevation level 2 ("card") — the tier used for floated
// surfaces like the login/signup card. Flat (hairline border, no shadow)
// is the baseline for most cards.
export const singleShadowTier = "rgba(15, 15, 15, 0.08) 0px 4px 12px 0px";

// Shared card corner radius (DESIGN.md `{rounded.lg}`) for pages that
// build their own Paper surfaces outside the theme's MuiPaper override.
export const cardRadius = radius.lg;

export const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: colors.primary,
      dark: colors.primaryPressed,
      light: colors.hairline,
      contrastText: colors.onPrimary,
    },
    info: {
      main: colors.linkBlue,
      dark: colors.linkBluePressed,
    },
    success: {
      main: colors.success,
    },
    warning: {
      main: colors.warning,
    },
    error: {
      main: colors.error,
    },
    text: {
      primary: colors.ink,
      secondary: colors.slate,
      disabled: colors.muted,
    },
    background: {
      // Subtle off-white app canvas with white card surfaces on top —
      // DESIGN.md's `{colors.surface}` / `{colors.canvas}` pairing.
      default: colors.surface,
      paper: colors.canvas,
    },
    divider: colors.hairline,
  },
  shape: {
    borderRadius: radius.md,
  },
  typography: {
    fontFamily,
    // DESIGN.md `typography` tokens, mapped onto the closest MUI variant.
    h1: { fontSize: 48, fontWeight: 600, lineHeight: 1.15, letterSpacing: "-0.5px" }, // heading-1
    h2: { fontSize: 36, fontWeight: 600, lineHeight: 1.2, letterSpacing: "-0.5px" }, // heading-2
    h3: { fontSize: 28, fontWeight: 600, lineHeight: 1.25, letterSpacing: 0 }, // heading-3
    h4: { fontSize: 22, fontWeight: 600, lineHeight: 1.3, letterSpacing: 0 }, // heading-4
    h5: { fontSize: 18, fontWeight: 600, lineHeight: 1.4, letterSpacing: 0 }, // heading-5
    h6: { fontSize: 18, fontWeight: 400, lineHeight: 1.5, letterSpacing: 0 }, // subtitle
    body1: { fontSize: 16, fontWeight: 400, lineHeight: 1.55, letterSpacing: 0 }, // body-md
    body2: { fontSize: 14, fontWeight: 400, lineHeight: 1.5, letterSpacing: 0 }, // body-sm
    // DESIGN.md only defines a bold caption (caption-bold); this is the
    // same size/line-height at regular weight for general helper text.
    caption: { fontSize: 13, fontWeight: 400, lineHeight: 1.4, letterSpacing: 0 },
    button: {
      fontSize: 14,
      fontWeight: 500,
      lineHeight: 1.3,
      letterSpacing: 0,
      textTransform: "none", // DESIGN.md buttons are sentence case, never all-caps
    },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: { backgroundColor: colors.surface, color: colors.ink },
      },
    },
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: ({ ownerState }) => ({
          borderRadius: radius.md,
          padding: "10px 18px",
          boxShadow: "none",
          ...(ownerState.variant === "contained" &&
            ownerState.color === "primary" && {
              backgroundColor: colors.primary,
              color: colors.onPrimary,
              "&:hover": { backgroundColor: colors.primaryPressed, boxShadow: "none" },
              "&:active": { backgroundColor: colors.primaryPressed },
              "&.Mui-disabled": { backgroundColor: colors.hairline, color: colors.muted },
            }),
        }),
        outlined: {
          borderColor: colors.hairlineStrong,
          color: colors.ink,
          "&:hover": { borderColor: colors.steel, backgroundColor: colors.surface },
        },
        text: {
          color: colors.ink,
          borderRadius: radius.sm,
          padding: "8px 12px",
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: radius.md,
          backgroundColor: colors.canvas,
          "& .MuiOutlinedInput-notchedOutline": { borderColor: colors.hairlineStrong },
          "&:hover .MuiOutlinedInput-notchedOutline": { borderColor: colors.steel },
          "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
            borderColor: colors.primary,
            borderWidth: 2,
          },
          "&.Mui-error .MuiOutlinedInput-notchedOutline": { borderColor: colors.error },
        },
        input: { padding: "14px 12px" },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: { color: colors.slate, "&.Mui-focused": { color: colors.primary } },
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
        rounded: { borderRadius: radius.lg },
      },
    },
    MuiAlert: {
      styleOverrides: {
        root: { borderRadius: radius.md },
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
    MuiToolbar: {
      styleOverrides: {
        // DESIGN.md Top Navigation is a ~64px sticky white bar.
        root: { minHeight: 64 },
      },
    },
    MuiLink: {
      styleOverrides: {
        // DESIGN.md reserves `{colors.primary}` for the dominant CTA only;
        // inline links use `{colors.link-blue}` — the two are never mixed.
        root: { color: colors.linkBlue, fontWeight: 500 },
      },
    },
  },
});
