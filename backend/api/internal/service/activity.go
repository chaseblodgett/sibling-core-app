package service

import (
	"context"
	"log"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

const (
	defaultActivityLimit = 20
	maxActivityLimit     = 100
)

// ActivityService reads a family's activity feed. Writing to it is not
// part of this interface — see logActivity, called directly by the other
// services (Task, Expense, Document, Parent, Family) as a side effect of
// their own writes, the same way those services already call into each
// other's repositories for cross-domain validation (requireMemberInFamily
// and friends).
type ActivityService interface {
	ListActivity(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error)
}

type activityService struct {
	activityRepo repository.ActivityRepository
}

// NewActivityService constructs an ActivityService.
func NewActivityService(activityRepo repository.ActivityRepository) ActivityService {
	return &activityService{activityRepo: activityRepo}
}

func (s *activityService) ListActivity(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
	if limit <= 0 {
		limit = defaultActivityLimit
	}
	if limit > maxActivityLimit {
		limit = maxActivityLimit
	}
	if offset < 0 {
		offset = 0
	}
	return s.activityRepo.ListByFamily(ctx, familyID, limit, offset)
}

// logActivity records one activity-feed entry, best-effort: a failure to
// log never fails the write that triggered it (same precedent as
// sendInviteEmail in family_service.go) — losing a feed entry is far
// cheaper than failing the task/expense/document/parent write the user
// is actually waiting on.
func logActivity(ctx context.Context, activityRepo repository.ActivityRepository, familyID string, actorID *string, actorName, verb, subjectType string, subjectID *string, subjectLabel string) {
	entry := &models.ActivityEntry{
		FamilyID:     familyID,
		ActorID:      actorID,
		ActorName:    actorName,
		Verb:         verb,
		SubjectType:  subjectType,
		SubjectID:    subjectID,
		SubjectLabel: subjectLabel,
	}
	if err := activityRepo.Create(ctx, entry); err != nil {
		log.Printf("activity log: failed to record %s on %s: %v", verb, subjectType, err)
	}
}
