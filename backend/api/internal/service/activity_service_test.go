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

func TestListActivityDefaultsLimit(t *testing.T) {
	repo := &mockrepository.MockActivityRepository{}
	svc := service.NewActivityService(repo)

	var gotLimit, gotOffset int
	repo.ListByFamilyFunc = func(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
		gotLimit, gotOffset = limit, offset
		return nil, nil
	}

	_, err := svc.ListActivity(context.Background(), "family-1", 0, 0)

	require.NoError(t, err)
	assert.Equal(t, 20, gotLimit)
	assert.Equal(t, 0, gotOffset)
}

func TestListActivityClampsOversizedLimit(t *testing.T) {
	repo := &mockrepository.MockActivityRepository{}
	svc := service.NewActivityService(repo)

	var gotLimit int
	repo.ListByFamilyFunc = func(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
		gotLimit = limit
		return nil, nil
	}

	_, err := svc.ListActivity(context.Background(), "family-1", 5000, 0)

	require.NoError(t, err)
	assert.Equal(t, 100, gotLimit)
}

func TestListActivityNegativeOffsetClampsToZero(t *testing.T) {
	repo := &mockrepository.MockActivityRepository{}
	svc := service.NewActivityService(repo)

	var gotOffset int
	repo.ListByFamilyFunc = func(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
		gotOffset = offset
		return nil, nil
	}

	_, err := svc.ListActivity(context.Background(), "family-1", 20, -5)

	require.NoError(t, err)
	assert.Equal(t, 0, gotOffset)
}
