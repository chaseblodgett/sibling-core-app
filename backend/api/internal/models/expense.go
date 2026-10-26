package models

import "time"

// Expense represents a cost incurred by a family member that may be split
// among other members of the family.
type Expense struct {
	ID          string    `json:"id" db:"id"`
	FamilyID    string    `json:"family_id" db:"family_id"`
	PaidBy      string    `json:"paid_by" db:"paid_by"` // member id who fronted the cost
	Description string    `json:"description" db:"description"`
	AmountCents int64     `json:"amount_cents" db:"amount_cents"`
	Currency    string    `json:"currency" db:"currency"`
	Category    string    `json:"category" db:"category"`
	ParentID    *string   `json:"parent_id,omitempty" db:"parent_id"`
	IncurredAt  time.Time `json:"incurred_at" db:"incurred_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Balance is a computed (not persisted) summary of what one member owes or
// is owed within a family, derived from expenses and their splits.
type Balance struct {
	MemberID string `json:"member_id"`
	NetCents int64  `json:"net_cents"` // positive = owed to them, negative = they owe
}

// SplitRequest carries one member's share in a create-expense request.
type SplitRequest struct {
	MemberID   string `json:"member_id"`
	ShareCents int64  `json:"share_cents"`
}

// CreateExpenseRequest is the JSON body for POST /families/:familyId/expenses.
// PaidBy defaults to the requester's own member id when omitted — most
// expenses are logged by whoever actually paid — but can name a different
// member (e.g. logging a sibling's receipt on their behalf).
type CreateExpenseRequest struct {
	Description string         `json:"description"`
	AmountCents int64          `json:"amount_cents"`
	Currency    string         `json:"currency"`
	Category    string         `json:"category"`
	ParentID    *string        `json:"parent_id,omitempty"`
	PaidBy      *string        `json:"paid_by,omitempty"`
	Splits      []SplitRequest `json:"splits"`
}

// SplitInput carries one member's share when logging an expense.
type SplitInput struct {
	MemberID   string
	ShareCents int64
}

// CreateExpenseInput carries the fields needed to log a new expense and
// how it's split. PaidBy and each Splits[i].MemberID are member ids
// (expenses.paid_by/expense_splits.member_id reference members, not
// users) — the handler must resolve ids for the family, same as Task.
type CreateExpenseInput struct {
	FamilyID    string
	PaidBy      string
	Description string
	AmountCents int64
	Currency    string
	Category    string
	ParentID    *string
	Splits      []SplitInput
}
