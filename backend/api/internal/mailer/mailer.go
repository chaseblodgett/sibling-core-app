// Package mailer sends transactional emails via Resend. Callers depend on
// the Mailer interface, not on Resend directly, so sending can be swapped
// or stubbed without touching business logic.
package mailer

import "context"

// FamilyInviteEmail carries everything needed to render an invite email.
type FamilyInviteEmail struct {
	To          string
	FamilyName  string
	InviterName string
	InviteCode  string
}

// Mailer sends transactional emails for the app.
type Mailer interface {
	SendFamilyInvite(ctx context.Context, email FamilyInviteEmail) error
}
