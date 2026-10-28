-- Per-user notification preferences. One row per user, created lazily on
-- first write (not pre-created for every signup) — a missing row means
-- "use the defaults", same semantics as a feature flag. digest_cadence
-- is reserved for a future digest-email phase — that phase needs a
-- scheduler subsystem this app doesn't have yet, so it's stored now
-- (cheap to add alongside the column this phase does use) but unused.

create table notification_preferences (
    user_id             uuid primary key references users(id) on delete cascade,
    task_assigned_email boolean not null default true,
    digest_cadence      text not null default 'none' check (digest_cadence in ('none', 'daily', 'weekly')),
    updated_at          timestamptz not null default now()
);
