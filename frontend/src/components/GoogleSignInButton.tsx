import { Alert } from "@mui/material";
import { useEffect, useRef, useState } from "react";

const CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID as string | undefined;

interface GoogleSignInButtonProps {
  text: "signin_with" | "signup_with";
  onIdToken: (idToken: string) => void;
}

// Renders Google's own "Sign in with Google" button via the Identity
// Services script (loaded in index.html). The script tag is async/defer,
// so this polls briefly for window.google to exist before initializing —
// it's almost always already there by the time React mounts, but can't be
// assumed.
export function GoogleSignInButton({ text, onIdToken }: GoogleSignInButtonProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [unavailable, setUnavailable] = useState(false);

  useEffect(() => {
    if (!CLIENT_ID) {
      setUnavailable(true);
      return;
    }

    let cancelled = false;
    const start = Date.now();

    const tryInit = () => {
      if (cancelled) return;
      if (!window.google) {
        if (Date.now() - start > 5000) {
          setUnavailable(true);
          return;
        }
        setTimeout(tryInit, 100);
        return;
      }

      window.google.accounts.id.initialize({
        client_id: CLIENT_ID,
        callback: (response) => onIdToken(response.credential),
      });
      if (containerRef.current) {
        window.google.accounts.id.renderButton(containerRef.current, {
          theme: "outline",
          size: "large",
          text,
          width: 350,
        });
      }
    };

    tryInit();
    return () => {
      cancelled = true;
    };
    // onIdToken is passed fresh each render from the caller but identifies
    // the same callback semantically; re-initializing on every change would
    // re-render Google's button and is unnecessary.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [text]);

  if (unavailable) {
    return <Alert severity="warning">Google sign-in isn't available right now.</Alert>;
  }

  return <div ref={containerRef} />;
}
