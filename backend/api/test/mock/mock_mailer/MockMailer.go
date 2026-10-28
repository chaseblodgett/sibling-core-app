package mockmailer

import (
	"context"

	"sibling-core/api/internal/mailer"
)

// MockMailer is a test double for mailer.Mailer. Assign SendFamilyInviteFunc
// in tests that care what was sent; unset it to no-op.
type MockMailer struct {
	SendFamilyInviteFunc func(ctx context.Context, email mailer.FamilyInviteEmail) error
	SendTaskAssignedFunc func(ctx context.Context, email mailer.TaskAssignedEmail) error
}

func (m *MockMailer) SendFamilyInvite(ctx context.Context, email mailer.FamilyInviteEmail) error {
	if m.SendFamilyInviteFunc == nil {
		return nil
	}
	return m.SendFamilyInviteFunc(ctx, email)
}

func (m *MockMailer) SendTaskAssigned(ctx context.Context, email mailer.TaskAssignedEmail) error {
	if m.SendTaskAssignedFunc == nil {
		return nil
	}
	return m.SendTaskAssignedFunc(ctx, email)
}
