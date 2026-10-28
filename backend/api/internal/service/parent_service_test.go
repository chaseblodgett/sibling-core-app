package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func newTestParentService(t *testing.T) (service.ParentService, *mockrepository.MockParentRepository) {
	t.Helper()
	repo := &mockrepository.MockParentRepository{}
	return service.NewParentService(repo), repo
}

func TestCreateParentSuccess(t *testing.T) {
	svc, repo := newTestParentService(t)

	var created *models.Parent
	repo.CreateFunc = func(ctx context.Context, p *models.Parent) error {
		p.ID = "parent-1"
		created = p
		return nil
	}

	parent, err := svc.CreateParent(context.Background(), models.CreateParentInput{
		FamilyID:  "family-1",
		Name:      "  Linda Doe ",
		Relation:  " Mom ",
		CreatedBy: "user-1",
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "Linda Doe", parent.Name)
	assert.Equal(t, "Mom", parent.Relation)
	assert.Equal(t, "family-1", parent.FamilyID)
}

func TestCreateParentBlankName(t *testing.T) {
	svc, _ := newTestParentService(t)

	_, err := svc.CreateParent(context.Background(), models.CreateParentInput{FamilyID: "family-1", Name: "   "})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUpdateParentPartial(t *testing.T) {
	svc, repo := newTestParentService(t)

	repo.GetByIDFunc = func(ctx context.Context, id string) (*models.Parent, error) {
		return &models.Parent{ID: id, Name: "Linda Doe", Relation: "Mom"}, nil
	}
	var updated *models.Parent
	repo.UpdateFunc = func(ctx context.Context, p *models.Parent) error {
		updated = p
		return nil
	}

	newRelation := "Stepmom"
	parent, err := svc.UpdateParent(context.Background(), "parent-1", models.UpdateParentInput{Relation: &newRelation})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Linda Doe", parent.Name, "name should be unchanged when not provided")
	assert.Equal(t, "Stepmom", parent.Relation)
}

func TestUpdateParentKeyInfo(t *testing.T) {
	svc, repo := newTestParentService(t)

	repo.GetByIDFunc = func(ctx context.Context, id string) (*models.Parent, error) {
		return &models.Parent{ID: id, Name: "Linda Doe", Relation: "Mom"}, nil
	}
	var updated *models.Parent
	repo.UpdateFunc = func(ctx context.Context, p *models.Parent) error {
		updated = p
		return nil
	}

	doctorName := "  Dr. Patel "
	allergies := " Penicillin "
	emergencyPhone := " 555-0100 "
	parent, err := svc.UpdateParent(context.Background(), "parent-1", models.UpdateParentInput{
		DoctorName:            &doctorName,
		Allergies:             &allergies,
		EmergencyContactPhone: &emergencyPhone,
	})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Dr. Patel", parent.DoctorName)
	assert.Equal(t, "Penicillin", parent.Allergies)
	assert.Equal(t, "555-0100", parent.EmergencyContactPhone)
	assert.Equal(t, "", parent.Medications, "fields not provided should stay unchanged")
	assert.Equal(t, "Linda Doe", parent.Name, "name should be unchanged when not provided")
}

func TestUpdateParentRejectsBlankName(t *testing.T) {
	svc, repo := newTestParentService(t)

	repo.GetByIDFunc = func(ctx context.Context, id string) (*models.Parent, error) {
		return &models.Parent{ID: id, Name: "Linda Doe"}, nil
	}

	blank := "   "
	_, err := svc.UpdateParent(context.Background(), "parent-1", models.UpdateParentInput{Name: &blank})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestListParents(t *testing.T) {
	svc, repo := newTestParentService(t)

	repo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Parent, error) {
		assert.Equal(t, "family-1", familyID)
		return []*models.Parent{{ID: "parent-1"}, {ID: "parent-2"}}, nil
	}

	parents, err := svc.ListParents(context.Background(), "family-1")

	require.NoError(t, err)
	assert.Len(t, parents, 2)
}
