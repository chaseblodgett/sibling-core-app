-- Read-only viewer role: an out-of-town relative or in-law who should
-- see what's going on without being able to change anything. Assigned
-- by an owner changing an existing member's role after they've already
-- joined normally (the shared family invite_code has no natural way to
-- carry a role, since it's one code used by everyone who joins).

alter table members drop constraint members_role_check;
alter table members add constraint members_role_check check (role in ('owner', 'member', 'viewer'));
