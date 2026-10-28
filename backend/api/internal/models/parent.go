package models

import "time"

// Parent represents a care-recipient tracked by a family — e.g. "Mom" or
// "Dad". Not a login/member: no auth, no credentials. Everyone in the
// family can see and manage a family's parents.
//
// The key-info fields (DoctorName through EmergencyContactPhone) are
// freeform text, not a structured medical record — the goal is "a
// sibling can read this fast in an emergency," not a drug database.
// They're edit-only (not settable on create) to keep the quick "add a
// parent" flow to just name + relation.
type Parent struct {
	ID                    string    `json:"id" db:"id"`
	FamilyID              string    `json:"family_id" db:"family_id"`
	Name                  string    `json:"name" db:"name"`
	Relation              string    `json:"relation" db:"relation"`
	DoctorName            string    `json:"doctor_name" db:"doctor_name"`
	DoctorPhone           string    `json:"doctor_phone" db:"doctor_phone"`
	Allergies             string    `json:"allergies" db:"allergies"`
	Medications           string    `json:"medications" db:"medications"`
	InsuranceProvider     string    `json:"insurance_provider" db:"insurance_provider"`
	InsuranceMemberID     string    `json:"insurance_member_id" db:"insurance_member_id"`
	EmergencyContactName  string    `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone" db:"emergency_contact_phone"`
	CreatedBy             string    `json:"created_by" db:"created_by"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// CreateParentRequest is the JSON body for POST /families/:familyId/parents.
type CreateParentRequest struct {
	Name     string `json:"name"`
	Relation string `json:"relation"`
}

// UpdateParentRequest is the JSON body for PATCH /parents/:id.
type UpdateParentRequest struct {
	Name                  *string `json:"name,omitempty"`
	Relation              *string `json:"relation,omitempty"`
	DoctorName            *string `json:"doctor_name,omitempty"`
	DoctorPhone           *string `json:"doctor_phone,omitempty"`
	Allergies             *string `json:"allergies,omitempty"`
	Medications           *string `json:"medications,omitempty"`
	InsuranceProvider     *string `json:"insurance_provider,omitempty"`
	InsuranceMemberID     *string `json:"insurance_member_id,omitempty"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
}

// CreateParentInput carries the fields needed to create a parent.
type CreateParentInput struct {
	FamilyID  string
	Name      string
	Relation  string
	CreatedBy string
	// ActorID/ActorName identify who's adding this parent, for the
	// activity-feed entry this create produces. CreatedBy above is a
	// user id (not a member id, unlike Task/Expense/Document), so it
	// can't double as ActorID the way CreateTaskInput.CreatedBy does.
	ActorID   string
	ActorName string
}

// UpdateParentInput carries the fields that may be updated on a parent.
type UpdateParentInput struct {
	Name                  *string
	Relation              *string
	DoctorName            *string
	DoctorPhone           *string
	Allergies             *string
	Medications           *string
	InsuranceProvider     *string
	InsuranceMemberID     *string
	EmergencyContactName  *string
	EmergencyContactPhone *string
}
