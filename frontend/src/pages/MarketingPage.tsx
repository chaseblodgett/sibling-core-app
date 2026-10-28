import AssignmentIcon from "@mui/icons-material/Assignment";
import AttachMoneyIcon from "@mui/icons-material/AttachMoney";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";
import NotificationsActiveIcon from "@mui/icons-material/NotificationsActive";
import {
  AppBar,
  Box,
  Button,
  Chip,
  Container,
  Grid,
  Paper,
  Stack,
  Toolbar,
  Typography,
} from "@mui/material";
import { Link as RouterLink } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { cardRadius } from "../theme/theme";

// Straight from utils/familyColor.ts's PALETTE (DESIGN.md's card-tint-*
// tokens) — the marketing page borrows the app's own family-color
// language for its scenario cards and hero preview, rather than
// inventing a separate marketing palette.
const peach = { bg: "#ffe8d4", border: "#dd5b00", text: "#793400" };
const mint = { bg: "#d9f3e1", border: "#1aae39", text: "#1aae39" };
const sky = { bg: "#dcecfa", border: "#0075de", text: "#005bab" };
const lavender = { bg: "#e6e0f5", border: "#7b3ff2", text: "#391c57" };

const MOMENTS = [
  {
    icon: <AssignmentIcon fontSize="small" />,
    color: peach,
    quote: "“Who's taking Dad to his 2 o'clock?”",
    body: "Tasks get assigned, not just asked about in a thread nobody scrolls back through. Everyone can see what's open, what's done, and whose turn it is.",
  },
  {
    icon: <AttachMoneyIcon fontSize="small" />,
    color: mint,
    quote: "“I covered the copay — did you see my message?”",
    body: "Log an expense once and Caretree splits it, tracks who's settled up, and keeps the receipt attached. No more reconstructing who owes what from memory.",
  },
  {
    icon: <InsertDriveFileIcon fontSize="small" />,
    color: sky,
    quote: "“Does anyone know where the insurance card is?”",
    body: "Upload it once. Every sibling can find it later — the POA, the insurance card, the discharge papers — without digging through old texts.",
  },
  {
    icon: <NotificationsActiveIcon fontSize="small" />,
    color: lavender,
    quote: "“What happened while I was at work?”",
    body: "A simple feed of what's been done, and by whom, so catching up doesn't mean re-asking the whole family what you missed.",
  },
];

const STEPS = [
  {
    title: "Create your family",
    body: "Start a family in Caretree and you'll get a short invite code to share.",
  },
  {
    title: "Invite your siblings",
    body: "Everyone joins with the code in seconds — no new account setup to walk anyone through.",
  },
  {
    title: "Coordinate for real",
    body: "Tasks, expenses, and documents live in one shared place instead of three different group chats.",
  },
];

export function MarketingPage() {
  const { user } = useAuth();
  const primaryCta = user ? { label: "Go to Dashboard", to: "/app" } : { label: "Sign up free", to: "/signup" };

  return (
    <Box sx={{ bgcolor: "background.default", minHeight: "100vh" }}>
      <AppBar position="static" color="transparent" elevation={0}>
        <Toolbar sx={{ maxWidth: 1200, width: "100%", mx: "auto" }}>
          <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 700, color: "text.primary" }}>
            Caretree
          </Typography>
          <Stack direction="row" spacing={1.5} sx={{ alignItems: "center" }}>
            {!user && (
              <Button component={RouterLink} to="/login" color="inherit">
                Log in
              </Button>
            )}
            <Button component={RouterLink} to={primaryCta.to} variant="contained">
              {primaryCta.label}
            </Button>
          </Stack>
        </Toolbar>
      </AppBar>

      <Container sx={{ maxWidth: 1200, pt: { xs: 6, md: 10 }, pb: { xs: 8, md: 12 } }}>
        <Grid container spacing={6} sx={{ alignItems: "center" }}>
          <Grid size={{ xs: 12, md: 6 }}>
            <Typography
              variant="h1"
              sx={{ fontSize: { xs: 36, sm: 44, md: 48 }, lineHeight: 1.1, mb: 2.5 }}
            >
              The group chat isn't working.
            </Typography>
            <Typography variant="body1" color="text.secondary" sx={{ mb: 4, maxWidth: 440, fontSize: 18 }}>
              Caretree gives siblings one shared place to track tasks, split expenses, and keep
              documents — for taking care of a parent together, without the guesswork over who's
              doing what.
            </Typography>
            <Stack direction="row" spacing={2}>
              <Button component={RouterLink} to={primaryCta.to} variant="contained" size="large">
                {primaryCta.label}
              </Button>
              {!user && (
                <Button component={RouterLink} to="/login" variant="outlined" size="large">
                  Log in
                </Button>
              )}
            </Stack>
          </Grid>

          <Grid size={{ xs: 12, md: 6 }}>
            <Paper
              sx={{
                borderRadius: `${cardRadius}px`,
                overflow: "hidden",
                border: "1px solid",
                borderColor: "divider",
                boxShadow: "rgba(15, 15, 15, 0.12) 0px 16px 40px -12px",
              }}
            >
              <Box sx={{ bgcolor: peach.bg, px: 2.5, py: 1.5 }}>
                <Typography variant="subtitle1" sx={{ fontWeight: 600, color: peach.text }}>
                  The Garcias
                </Typography>
              </Box>
              <Stack divider={<Box sx={{ borderBottom: 1, borderColor: "divider" }} />}>
                <Stack direction="row" spacing={1} sx={{ alignItems: "center", px: 2.5, py: 1.75 }}>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="body2">Pick up prescription</Typography>
                  </Box>
                  <Chip label="Oct 24" size="small" color="error" variant="outlined" />
                </Stack>
                <Stack direction="row" spacing={1} sx={{ alignItems: "center", px: 2.5, py: 1.75 }}>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="body2" sx={{ textDecoration: "line-through", color: "text.disabled" }}>
                      Call about insurance
                    </Typography>
                  </Box>
                  <Chip label="Done" size="small" color="success" variant="outlined" />
                </Stack>
                <Stack direction="row" spacing={1} sx={{ alignItems: "center", px: 2.5, py: 1.75 }}>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="body2">Grandma's birthday dinner</Typography>
                  </Box>
                  <Chip label="Jamie" size="small" variant="outlined" />
                </Stack>
              </Stack>
            </Paper>
          </Grid>
        </Grid>
      </Container>

      <Box sx={{ bgcolor: "background.paper", py: { xs: 8, md: 10 } }}>
        <Container sx={{ maxWidth: 1200 }}>
          <Typography variant="h2" sx={{ fontSize: { xs: 28, md: 36 }, mb: 1 }}>
            The questions every caregiving family asks
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 5, maxWidth: 560 }}>
            Caretree answers them before anyone has to ask.
          </Typography>
          <Grid container spacing={3}>
            {MOMENTS.map((moment) => (
              <Grid key={moment.quote} size={{ xs: 12, sm: 6 }}>
                <Paper
                  sx={{
                    bgcolor: moment.color.bg,
                    borderRadius: `${cardRadius}px`,
                    p: 3.5,
                    height: "100%",
                  }}
                >
                  <Box sx={{ color: moment.color.text, mb: 1.5 }}>{moment.icon}</Box>
                  <Typography variant="h6" component="p" sx={{ color: moment.color.text, mb: 1.5 }}>
                    {moment.quote}
                  </Typography>
                  <Typography variant="body2" sx={{ color: moment.color.text, opacity: 0.85 }}>
                    {moment.body}
                  </Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      <Container sx={{ maxWidth: 1200, py: { xs: 8, md: 10 } }}>
        <Typography variant="h2" sx={{ fontSize: { xs: 28, md: 36 }, mb: 5, textAlign: "center" }}>
          Getting started takes minutes
        </Typography>
        <Grid container spacing={4}>
          {STEPS.map((step, i) => (
            <Grid key={step.title} size={{ xs: 12, sm: 4 }}>
              <Stack spacing={1.5}>
                <Box
                  sx={{
                    width: 36,
                    height: 36,
                    borderRadius: "50%",
                    bgcolor: "primary.main",
                    color: "primary.contrastText",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontWeight: 600,
                  }}
                >
                  {i + 1}
                </Box>
                <Typography variant="h6" component="h3">
                  {step.title}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {step.body}
                </Typography>
              </Stack>
            </Grid>
          ))}
        </Grid>
      </Container>

      <Box sx={{ bgcolor: "primary.main", py: { xs: 7, md: 9 } }}>
        <Container sx={{ maxWidth: 1200, textAlign: "center" }}>
          <Typography variant="h2" sx={{ fontSize: { xs: 28, md: 36 }, color: "primary.contrastText", mb: 1.5 }}>
            Start coordinating today
          </Typography>
          <Typography variant="body1" sx={{ color: "primary.contrastText", opacity: 0.85, mb: 4 }}>
            Free to start. Invite your siblings whenever you're ready.
          </Typography>
          <Button
            component={RouterLink}
            to={primaryCta.to}
            variant="contained"
            size="large"
            sx={{ bgcolor: "common.white", color: "primary.main", "&:hover": { bgcolor: "grey.100" } }}
          >
            {primaryCta.label}
          </Button>
        </Container>
      </Box>

      <Box component="footer" sx={{ py: 4, borderTop: 1, borderColor: "divider" }}>
        <Container sx={{ maxWidth: 1200 }}>
          <Typography variant="body2" color="text.secondary">
            Caretree — built for siblings coordinating care for a parent.
          </Typography>
        </Container>
      </Box>
    </Box>
  );
}
