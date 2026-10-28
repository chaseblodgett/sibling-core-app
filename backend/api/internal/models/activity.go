package models

import "time"

// ActivityEntry is one row in a family's activity feed. It's a
// self-contained snapshot of what happened — ActorName and SubjectLabel
// are captured at write time rather than resolved via a join at read
// time — so the feed keeps reading correctly even after the member who
// did something, or the task/parent/document it concerned, is later
// removed or deleted. ActorID/SubjectID are kept for future aggregation
// (e.g. a contribution leaderboard as a query over this same table) but
// aren't required for rendering the feed itself.
type ActivityEntry struct {
	ID           string    `json:"id" db:"id"`
	FamilyID     string    `json:"family_id" db:"family_id"`
	ActorID      *string   `json:"actor_id,omitempty" db:"actor_id"`
	ActorName    string    `json:"actor_name" db:"actor_name"`
	Verb         string    `json:"verb" db:"verb"`
	SubjectType  string    `json:"subject_type" db:"subject_type"`
	SubjectID    *string   `json:"subject_id,omitempty" db:"subject_id"`
	SubjectLabel string    `json:"subject_label" db:"subject_label"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
