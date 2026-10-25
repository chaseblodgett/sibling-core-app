-- Custom email/password auth: the Go backend owns user identity directly
-- (bcrypt password hashes, backend-issued JWTs) rather than delegating to
-- Supabase Auth. This repoints the FKs that previously referenced
-- Supabase's auth.users at this table instead.

create table users (
    id            uuid primary key default gen_random_uuid(),
    email         text not null unique,
    password_hash text not null,
    name          text not null,
    created_at    timestamptz not null default now()
);

alter table families
    drop constraint families_created_by_fkey,
    add constraint families_created_by_fkey
        foreign key (created_by) references users(id) on delete restrict;

alter table family_invites
    drop constraint family_invites_invited_by_fkey,
    add constraint family_invites_invited_by_fkey
        foreign key (invited_by) references users(id) on delete restrict;

alter table members
    drop constraint members_user_id_fkey,
    add constraint members_user_id_fkey
        foreign key (user_id) references users(id) on delete cascade;
