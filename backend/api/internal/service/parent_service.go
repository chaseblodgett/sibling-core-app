package service

import (
	"context"
	"fmt"
	"strings"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// ParentService contains business logic for a family's tracked parents
// (care-recipients). Family-membership authorization (who may call these
// at all, and who may delete) lives in the handler via FamilyService.
type ParentService interface {
	CreateParent(ctx context.Context, input models.CreateParentInput) (*models.Parent, error)
	GetParent(ctx context.Context, id string) (*models.Parent, error)
	ListParents(ctx context.Context, familyID string) ([]*models.Parent, error)
	UpdateParent(ctx context.Context, id string, input models.UpdateParentInput) (*models.Parent, error)
	DeleteParent(ctx context.Context, id string) error
}

type parentService struct {
	parentRepo repository.ParentRepository
}

// NewParentService constructs a ParentService.
func NewParentService(parentRepo repository.ParentRepository) ParentService {
	return &parentService{parentRepo: parentRepo}
}

func (s *parentService) CreateParent(ctx context.Context, input models.CreateParentInput) (*models.Parent, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	parent := &models.Parent{
		FamilyID:  input.FamilyID,
		Name:      name,
		Relation:  strings.TrimSpace(input.Relation),
		CreatedBy: input.CreatedBy,
	}
	if err := s.parentRepo.Create(ctx, parent); err != nil {
		return nil, err
	}
	return parent, nil
}

func (s *parentService) GetParent(ctx context.Context, id string) (*models.Parent, error) {
	return s.parentRepo.GetByID(ctx, id)
}

func (s *parentService) ListParents(ctx context.Context, familyID string) ([]*models.Parent, error) {
	return s.parentRepo.ListByFamily(ctx, familyID)
}

func (s *parentService) UpdateParent(ctx context.Context, id string, input models.UpdateParentInput) (*models.Parent, error) {
	parent, err := s.parentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
		}
		parent.Name = name
	}
	if input.Relation != nil {
		parent.Relation = strings.TrimSpace(*input.Relation)
	}
	if input.DoctorName != nil {
		parent.DoctorName = strings.TrimSpace(*input.DoctorName)
	}
	if input.DoctorPhone != nil {
		parent.DoctorPhone = strings.TrimSpace(*input.DoctorPhone)
	}
	if input.Allergies != nil {
		parent.Allergies = strings.TrimSpace(*input.Allergies)
	}
	if input.Medications != nil {
		parent.Medications = strings.TrimSpace(*input.Medications)
	}
	if input.InsuranceProvider != nil {
		parent.InsuranceProvider = strings.TrimSpace(*input.InsuranceProvider)
	}
	if input.InsuranceMemberID != nil {
		parent.InsuranceMemberID = strings.TrimSpace(*input.InsuranceMemberID)
	}
	if input.EmergencyContactName != nil {
		parent.EmergencyContactName = strings.TrimSpace(*input.EmergencyContactName)
	}
	if input.EmergencyContactPhone != nil {
		parent.EmergencyContactPhone = strings.TrimSpace(*input.EmergencyContactPhone)
	}

	if err := s.parentRepo.Update(ctx, parent); err != nil {
		return nil, err
	}
	return parent, nil
}

func (s *parentService) DeleteParent(ctx context.Context, id string) error {
	return s.parentRepo.Delete(ctx, id)
}
