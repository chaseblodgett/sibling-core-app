package models

// Contribution is a computed (not persisted) per-member summary of how
// much caregiving work a family member has picked up — tasks they've
// actually completed (not just been assigned) and expenses they've
// actually fronted (not just logged on someone else's behalf). It's a
// read-side aggregation over Task/Expense/ActivityEntry data, the same
// way Balance is a computed aggregation over Expense/ExpenseSplit data.
type Contribution struct {
	MemberID          string `json:"member_id"`
	MemberName        string `json:"member_name"`
	TasksCompleted    int    `json:"tasks_completed"`
	ExpensesPaidCount int    `json:"expenses_paid_count"`
	ExpensesPaidCents int64  `json:"expenses_paid_cents"`
}
