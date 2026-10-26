package models

import "time"

// Parent represents a care-recipient tracked by a family — e.g. "Mom" or
// "Dad". Not a login/member: no auth, no credentials. Everyone in the
// family can see and manage a family's parents.
type Parent struct {
	ID        string    `json:"id" db:"id"`
	FamilyID  string    `json:"family_id" db:"family_id"`
	Name      string    `json:"name" db:"name"`
	Relation  string    `json:"relation" db:"relation"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateParentRequest is the JSON body for POST /families/:familyId/parents.
type CreateParentRequest struct {
	Name     string `json:"name"`
	Relation string `json:"relation"`
}

// UpdateParentRequest is the JSON body for PATCH /parents/:id.
type UpdateParentRequest struct {
	Name     *string `json:"name,omitempty"`
	Relation *string `json:"relation,omitempty"`
}

// CreateParentInput carries the fields needed to create a parent.
type CreateParentInput struct {
	FamilyID  string
	Name      string
	Relation  string
	CreatedBy string
}

// UpdateParentInput carries the fields that may be updated on a parent.
type UpdateParentInput struct {
	Name     *string
	Relation *string
}
