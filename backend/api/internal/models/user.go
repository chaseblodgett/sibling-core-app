package models

import "time"

// User is an account in the app's own auth system — either email/password
// or "Sign in with Google". PasswordHash is nil for a Google-only account
// (one that has never set a password); GoogleSub is nil for a user who has
// never signed in with Google. Neither is ever serialized to JSON.
type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash *string   `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	GoogleSub    *string   `json:"-" db:"google_sub"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
