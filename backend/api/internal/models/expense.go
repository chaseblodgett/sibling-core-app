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
	IncurredAt  time.Time `json:"incurred_at" db:"incurred_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Balance is a computed (not persisted) summary of what one member owes or
// is owed within a family, derived from expenses and their splits.
type Balance struct {
	MemberID    string `json:"member_id"`
	NetCents    int64  `json:"net_cents"` // positive = owed to them, negative = they owe
}

// SplitRequest carries one member's share in a create-expense request.
type SplitRequest struct {
	MemberID   string `json:"member_id"`
	ShareCents int64  `json:"share_cents"`
}

// CreateExpenseRequest is the JSON body for POST /families/:familyId/expenses.
type CreateExpenseRequest struct {
	Description string         `json:"description"`
	AmountCents int64          `json:"amount_cents"`
	Currency    string         `json:"currency"`
	Category    string         `json:"category"`
	Splits      []SplitRequest `json:"splits"`
}

// SplitInput carries one member's share when logging an expense.
type SplitInput struct {
	MemberID   string
	ShareCents int64
}

// CreateExpenseInput carries the fields needed to log a new expense and how it's split.
type CreateExpenseInput struct {
	FamilyID    string
	PaidBy      string
	Description string
	AmountCents int64
	Currency    string
	Category    string
	Splits      []SplitInput
}
