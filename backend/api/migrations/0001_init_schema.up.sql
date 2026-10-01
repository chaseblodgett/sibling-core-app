-- Initial schema for families, membership, tasks, expenses, and documents.
-- Written for Supabase Postgres: members.user_id references the built-in
-- auth.users table managed by Supabase Auth.

create extension if not exists pgcrypto;

create table families (
    id          uuid primary key default gen_random_uuid(),
    name        text not null,
    invite_code text not null unique,
    created_by  uuid not null references auth.users(id) on delete restrict,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create table family_invites (
    id          uuid primary key default gen_random_uuid(),
    family_id   uuid not null references families(id) on delete cascade,
    email       text not null,
    invited_by  uuid not null references auth.users(id) on delete restrict,
    status      text not null default 'pending' check (status in ('pending', 'accepted', 'revoked')),
    created_at  timestamptz not null default now(),
    accepted_at timestamptz
);

create index idx_family_invites_family_id on family_invites(family_id);
create index idx_family_invites_email on family_invites(email);

create table members (
    id         uuid primary key default gen_random_uuid(),
    family_id  uuid not null references families(id) on delete cascade,
    user_id    uuid not null references auth.users(id) on delete cascade,
    email      text not null,
    full_name  text not null default '',
    role       text not null default 'member' check (role in ('owner', 'member')),
    joined_at  timestamptz not null default now(),
    unique (family_id, user_id)
);

create index idx_members_family_id on members(family_id);
create index idx_members_user_id on members(user_id);

create table tasks (
    id          uuid primary key default gen_random_uuid(),
    family_id   uuid not null references families(id) on delete cascade,
    title       text not null,
    description text not null default '',
    status      text not null default 'todo' check (status in ('todo', 'in_progress', 'done')),
    assignee_id uuid references members(id) on delete set null,
    created_by  uuid not null references members(id) on delete restrict,
    due_date    timestamptz,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create index idx_tasks_family_id on tasks(family_id);
create index idx_tasks_assignee_id on tasks(assignee_id);

create table expenses (
    id           uuid primary key default gen_random_uuid(),
    family_id    uuid not null references families(id) on delete cascade,
    paid_by      uuid not null references members(id) on delete restrict,
    description  text not null,
    amount_cents bigint not null check (amount_cents >= 0),
    currency     text not null default 'USD',
    category     text not null default '',
    incurred_at  timestamptz not null default now(),
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now()
);

create index idx_expenses_family_id on expenses(family_id);
create index idx_expenses_paid_by on expenses(paid_by);

create table expense_splits (
    id           uuid primary key default gen_random_uuid(),
    expense_id   uuid not null references expenses(id) on delete cascade,
    member_id    uuid not null references members(id) on delete cascade,
    share_cents  bigint not null check (share_cents >= 0),
    settled      boolean not null default false,
    unique (expense_id, member_id)
);

create index idx_expense_splits_expense_id on expense_splits(expense_id);
create index idx_expense_splits_member_id on expense_splits(member_id);

create table documents (
    id           uuid primary key default gen_random_uuid(),
    family_id    uuid not null references families(id) on delete cascade,
    uploaded_by  uuid not null references members(id) on delete restrict,
    file_name    text not null,
    storage_path text not null,
    content_type text not null default 'application/octet-stream',
    size_bytes   bigint not null default 0,
    created_at   timestamptz not null default now()
);

create index idx_documents_family_id on documents(family_id);
