import { createClient } from "@supabase/supabase-js";

const supabaseUrl = import.meta.env.VITE_SUPABASE_URL ?? "";
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY ?? "";

// Used for Supabase Auth (sign up/login sessions) and direct Storage access
// from the client. Application data still goes through the Go API.
export const supabase = createClient(supabaseUrl, supabaseAnonKey);
