package handler

import (
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// memberDisplayName returns a member's display name for activity-feed
// attribution — full name when set, email otherwise (same fallback the
// frontend already uses for members without a name).
func memberDisplayName(m *models.Member) string {
	if m.FullName != "" {
		return m.FullName
	}
	return m.Email
}

// requireNotViewer rejects a write from a viewer-role member. Every
// create/update/delete handler for a shared family resource calls this
// right after resolving the requester's membership — viewers can read
// everything but change nothing.
func requireNotViewer(member *models.Member) error {
	if member.Role == "viewer" {
		return service.ErrNotAuthorized
	}
	return nil
}
