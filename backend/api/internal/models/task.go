package models

import "time"

// Task represents a shared to-do item on a family's task board.
type Task struct {
	ID             string     `json:"id" db:"id"`
	FamilyID       string     `json:"family_id" db:"family_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Status         string     `json:"status" db:"status"` // todo, in_progress, done
	AssigneeID     *string    `json:"assignee_id,omitempty" db:"assignee_id"`
	ParentID       *string    `json:"parent_id,omitempty" db:"parent_id"`
	CreatedBy      string     `json:"created_by" db:"created_by"`
	DueDate        *time.Time `json:"due_date,omitempty" db:"due_date"`
	RecurrenceRule string     `json:"recurrence_rule" db:"recurrence_rule"` // none, daily, weekly, monthly
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateTaskRequest is the JSON body for POST /families/:familyId/tasks.
type CreateTaskRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	AssigneeID     *string    `json:"assignee_id,omitempty"`
	ParentID       *string    `json:"parent_id,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	RecurrenceRule string     `json:"recurrence_rule,omitempty"`
}

// UpdateTaskRequest is the JSON body for PATCH /tasks/:id.
type UpdateTaskRequest struct {
	Title          *string    `json:"title,omitempty"`
	Description    *string    `json:"description,omitempty"`
	Status         *string    `json:"status,omitempty"`
	AssigneeID     *string    `json:"assignee_id,omitempty"`
	ParentID       *string    `json:"parent_id,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	RecurrenceRule *string    `json:"recurrence_rule,omitempty"`
}

// AssignTaskRequest is the JSON body for POST /tasks/:id/assign.
type AssignTaskRequest struct {
	MemberID string `json:"member_id"`
}

// CreateTaskInput carries the fields needed to create a task. CreatedBy
// and AssigneeID are member ids (tasks.created_by/assignee_id reference
// members, not users) — the handler must resolve the requester's member
// row for the family first.
type CreateTaskInput struct {
	FamilyID       string
	Title          string
	Description    string
	AssigneeID     *string
	ParentID       *string
	CreatedBy      string
	DueDate        *time.Time
	RecurrenceRule string
	// ActorName is CreatedBy's display name, for the activity-feed entry
	// this create produces — snapshotted at write time rather than
	// resolved by a join later.
	ActorName string
}

// UpdateTaskInput carries the fields that may be updated on a task.
type UpdateTaskInput struct {
	Title          *string
	Description    *string
	Status         *string
	AssigneeID     *string
	ParentID       *string
	DueDate        *time.Time
	RecurrenceRule *string
	// ActorID/ActorName identify who's making this update, for the
	// activity-feed entry a completion produces. Unlike Create, Update
	// has no existing "who did this" field to reuse.
	ActorID   string
	ActorName string
}
