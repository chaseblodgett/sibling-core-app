package service

import (
	"context"
	"errors"
	"fmt"

	"sibling-core/api/internal/repository"
)

// requireMemberInFamily checks that memberID names a member of familyID,
// returning ErrInvalidInput (not a repository error) either way so
// callers can surface a clean 400 regardless of why it failed.
func requireMemberInFamily(ctx context.Context, memberRepo repository.MemberRepository, memberID, familyID, label string) error {
	member, err := memberRepo.GetByID(ctx, memberID)
	if err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return fmt.Errorf("%w: %s is not a member of this family", ErrInvalidInput, label)
		}
		return err
	}
	if member.FamilyID != familyID {
		return fmt.Errorf("%w: %s is not a member of this family", ErrInvalidInput, label)
	}
	return nil
}

// requireParentInFamily checks that parentID names a parent of familyID.
func requireParentInFamily(ctx context.Context, parentRepo repository.ParentRepository, parentID, familyID string) error {
	parent, err := parentRepo.GetByID(ctx, parentID)
	if err != nil {
		if errors.Is(err, repository.ErrParentNotFound) {
			return fmt.Errorf("%w: parent does not belong to this family", ErrInvalidInput)
		}
		return err
	}
	if parent.FamilyID != familyID {
		return fmt.Errorf("%w: parent does not belong to this family", ErrInvalidInput)
	}
	return nil
}

// requireDocumentInFamily checks that documentID names a document of familyID.
func requireDocumentInFamily(ctx context.Context, documentRepo repository.DocumentRepository, documentID, familyID string) error {
	doc, err := documentRepo.GetByID(ctx, documentID)
	if err != nil {
		if errors.Is(err, repository.ErrDocumentNotFound) {
			return fmt.Errorf("%w: document does not belong to this family", ErrInvalidInput)
		}
		return err
	}
	if doc.FamilyID != familyID {
		return fmt.Errorf("%w: document does not belong to this family", ErrInvalidInput)
	}
	return nil
}
