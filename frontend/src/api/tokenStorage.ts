// Wraps localStorage access for the auth token. See the auth implementation
// notes for why localStorage was chosen over an httpOnly cookie here.
const TOKEN_KEY = "sibling_core_token";

export const tokenStorage = {
  get(): string | null {
    try {
      return localStorage.getItem(TOKEN_KEY);
    } catch {
      return null;
    }
  },
  set(token: string): void {
    try {
      localStorage.setItem(TOKEN_KEY, token);
    } catch {
      // Storage unavailable (e.g. private browsing) — session just won't persist.
    }
  },
  clear(): void {
    try {
      localStorage.removeItem(TOKEN_KEY);
    } catch {
      // ignore
    }
  },
};
