package models

import "time"

// Family represents a household/family group that members belong to.
type Family struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	InviteCode string    `json:"invite_code" db:"invite_code"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// FamilyInvite represents a pending invitation for a user to join a family by email.
type FamilyInvite struct {
	ID         string     `json:"id" db:"id"`
	FamilyID   string     `json:"family_id" db:"family_id"`
	Email      string     `json:"email" db:"email"`
	InvitedBy  string     `json:"invited_by" db:"invited_by"`
	Status     string     `json:"status" db:"status"` // pending, accepted, revoked
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty" db:"accepted_at"`
}

// CreateFamilyRequest is the JSON body for POST /families.
type CreateFamilyRequest struct {
	Name string `json:"name"`
}

// JoinFamilyRequest is the JSON body for POST /families/join.
type JoinFamilyRequest struct {
	InviteCode string `json:"invite_code"`
}

// InviteMemberRequest is the JSON body for POST /families/:id/invites.
type InviteMemberRequest struct {
	Email string `json:"email"`
}

// UpdateMemberRoleRequest is the JSON body for
// PATCH /families/:id/members/:memberId/role.
type UpdateMemberRoleRequest struct {
	Role string `json:"role"`
}

// CreateFamilyInput carries the fields needed to create a new family.
type CreateFamilyInput struct {
	Name      string
	CreatedBy string // user id of the creator, who becomes the owner
}

// InviteMemberInput carries the fields needed to invite someone by email.
type InviteMemberInput struct {
	FamilyID  string
	Email     string
	InvitedBy string
}
