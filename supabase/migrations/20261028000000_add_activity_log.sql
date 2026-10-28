-- Activity feed: a reverse-chronological "what's been happening" list per
-- family, so siblings coordinating async can see what they missed instead
-- of only ever seeing current-state stats. Entries are self-contained
-- snapshots (actor_name, subject_label) rather than live joins — the feed
-- should keep reading correctly even after the member who did something,
-- or the task/parent/document it was about, is later removed or deleted.
-- actor_id/subject_id are kept for future aggregation (e.g. a
-- contribution leaderboard built as a query over this same table), but
-- subject_id deliberately has no foreign key: the subject may be deleted
-- without the activity entry disappearing with it.

create table activity_log (
    id            uuid primary key default gen_random_uuid(),
    family_id     uuid not null references families(id) on delete cascade,
    actor_id      uuid references members(id) on delete set null,
    actor_name    text not null,
    verb          text not null,
    subject_type  text not null,
    subject_id    uuid,
    subject_label text not null default '',
    created_at    timestamptz not null default now()
);

create index idx_activity_log_family_created on activity_log(family_id, created_at desc);
