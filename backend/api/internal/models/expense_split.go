package models

// ExpenseSplit represents one member's share of a given expense.
type ExpenseSplit struct {
	ID          string `json:"id" db:"id"`
	ExpenseID   string `json:"expense_id" db:"expense_id"`
	MemberID    string `json:"member_id" db:"member_id"`
	ShareCents  int64  `json:"share_cents" db:"share_cents"`
	Settled     bool   `json:"settled" db:"settled"`
}
