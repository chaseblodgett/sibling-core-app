-- Receipt attachment: Document and Expense were previously disconnected
-- — uploading a receipt photo meant dropping it in Documents and hoping
-- someone connects the dots. One receipt per expense to start (a join
-- table for multiple receipts can come later if real usage needs it).

alter table expenses add column document_id uuid references documents(id) on delete set null;
