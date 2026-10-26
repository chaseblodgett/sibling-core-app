package mailer

import (
	"context"
	"log"
)

// noopMailer is used when no RESEND_API_KEY is configured. Invites still
// get recorded; this just skips the send and says so, once, so local dev
// doesn't need a real email account to exercise the invite flow.
type noopMailer struct{}

// NewNoopMailer constructs a Mailer that logs instead of sending.
func NewNoopMailer() Mailer {
	return &noopMailer{}
}

func (m *noopMailer) SendFamilyInvite(_ context.Context, email FamilyInviteEmail) error {
	log.Printf("mailer: RESEND_API_KEY not set, skipping invite email to %s (family %q, code %s)",
		email.To, email.FamilyName, email.InviteCode)
	return nil
}
