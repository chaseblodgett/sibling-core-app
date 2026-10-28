-- Parent "key info": the things a sibling needs fast — stepping in
-- during an emergency, or on a call with a new doctor — instead of
-- hunting through uploaded documents or a group chat. Freeform text for
-- now (medications especially: a readable list, not a structured drug
-- database); add structure later only if real usage shows it's needed.

alter table parents
    add column doctor_name              text not null default '',
    add column doctor_phone             text not null default '',
    add column allergies                text not null default '',
    add column medications              text not null default '',
    add column insurance_provider       text not null default '',
    add column insurance_member_id      text not null default '',
    add column emergency_contact_name   text not null default '',
    add column emergency_contact_phone  text not null default '';
