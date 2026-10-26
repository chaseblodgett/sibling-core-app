-- Supports "Sign in with Google" accounts: password_hash is no longer
-- required (a Google-only account never sets one), and google_sub links a
-- user to their Google identity for repeat sign-ins and account linking.

alter table users
    alter column password_hash drop not null,
    add column google_sub text unique;
