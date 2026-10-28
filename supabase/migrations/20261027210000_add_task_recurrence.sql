-- Recurring tasks: med refills, weekly check-ins, monthly bill payments
-- are the normal shape of caregiving work. recurrence_rule describes how
-- a task repeats; completing a task whose rule isn't 'none' spawns the
-- next occurrence (handled in the service layer, not the database).

alter table tasks add column recurrence_rule text not null default 'none'
    check (recurrence_rule in ('none', 'daily', 'weekly', 'monthly'));
