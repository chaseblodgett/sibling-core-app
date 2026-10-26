package mockgoogleauth

import (
	"context"

	"sibling-core/api/internal/googleauth"
)

// MockVerifier is a test double for googleauth.Verifier. Assign VerifyFunc
// in tests that care about the claims returned; unset it to get a nil
// claims/nil error zero value.
type MockVerifier struct {
	VerifyFunc func(ctx context.Context, idToken string) (*googleauth.Claims, error)
}

func (m *MockVerifier) Verify(ctx context.Context, idToken string) (*googleauth.Claims, error) {
	if m.VerifyFunc == nil {
		return nil, nil
	}
	return m.VerifyFunc(ctx, idToken)
}
