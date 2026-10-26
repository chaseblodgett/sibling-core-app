-- Parents: care-recipients tracked by a family. Not a login/member — no
-- auth/credentials, just a name + relation. Tasks, expenses, and documents
-- can optionally tie back to a specific parent so a family can see what a
-- given parent needs; null parent_id means "general family item".

create table parents (
    id          uuid primary key default gen_random_uuid(),
    family_id   uuid not null references families(id) on delete cascade,
    name        text not null,
    relation    text not null default '',
    created_by  uuid not null references users(id) on delete restrict,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create index idx_parents_family_id on parents(family_id);

alter table tasks add column parent_id uuid references parents(id) on delete set null;
alter table expenses add column parent_id uuid references parents(id) on delete set null;
alter table documents add column parent_id uuid references parents(id) on delete set null;

create index idx_tasks_parent_id on tasks(parent_id);
create index idx_expenses_parent_id on expenses(parent_id);
create index idx_documents_parent_id on documents(parent_id);
