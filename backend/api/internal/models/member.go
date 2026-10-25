package models

import "time"

// Member represents a user's membership within a family group.
// The user identity itself (id, email, password) lives in the users table;
// this row links a user id to a family with a role.
type Member struct {
	ID       string    `json:"id" db:"id"`
	FamilyID string    `json:"family_id" db:"family_id"`
	UserID   string    `json:"user_id" db:"user_id"` // users.id
	Email    string    `json:"email" db:"email"`
	FullName string    `json:"full_name" db:"full_name"`
	Role     string    `json:"role" db:"role"` // owner, member
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}
