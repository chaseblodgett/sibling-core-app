package service_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
	mockstorage "sibling-core/api/test/mock/mock_storage"
)

func newTestDocumentService(t *testing.T) (service.DocumentService, *mockrepository.MockDocumentRepository, *mockrepository.MockParentRepository, *mockstorage.MockStorage) {
	t.Helper()
	documentRepo := &mockrepository.MockDocumentRepository{}
	parentRepo := &mockrepository.MockParentRepository{}
	store := &mockstorage.MockStorage{}
	return service.NewDocumentService(documentRepo, parentRepo, store), documentRepo, parentRepo, store
}

func TestUploadDocumentSuccess(t *testing.T) {
	svc, documentRepo, _, store := newTestDocumentService(t)

	var uploadedPath, uploadedContentType string
	store.UploadFunc = func(ctx context.Context, path, contentType string, content io.Reader) error {
		uploadedPath = path
		uploadedContentType = contentType
		body, err := io.ReadAll(content)
		require.NoError(t, err)
		assert.Equal(t, "file contents", string(body))
		return nil
	}

	var created *models.Document
	documentRepo.CreateFunc = func(ctx context.Context, d *models.Document) error {
		d.ID = "doc-1"
		created = d
		return nil
	}

	doc, err := svc.UploadDocument(context.Background(), models.UploadDocumentInput{
		FamilyID:    "family-1",
		UploadedBy:  "member-1",
		FileName:    "Mom's will (final).pdf",
		ContentType: "application/pdf",
		SizeBytes:   13,
		Content:     strings.NewReader("file contents"),
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "doc-1", doc.ID)
	assert.Equal(t, "Mom's will (final).pdf", doc.FileName)
	assert.Equal(t, "application/pdf", uploadedContentType)
	assert.Contains(t, uploadedPath, "family-1/")
	assert.NotContains(t, uploadedPath, "(")
	assert.NotContains(t, uploadedPath, " ")
}

func TestUploadDocumentRejectsEmptyFile(t *testing.T) {
	svc, _, _, _ := newTestDocumentService(t)

	_, err := svc.UploadDocument(context.Background(), models.UploadDocumentInput{
		FamilyID: "family-1", UploadedBy: "member-1", FileName: "empty.txt", SizeBytes: 0,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUploadDocumentRejectsBlankFileName(t *testing.T) {
	svc, _, _, _ := newTestDocumentService(t)

	_, err := svc.UploadDocument(context.Background(), models.UploadDocumentInput{
		FamilyID: "family-1", UploadedBy: "member-1", FileName: "  ", SizeBytes: 10,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUploadDocumentDefaultsContentType(t *testing.T) {
	svc, documentRepo, _, store := newTestDocumentService(t)

	var gotContentType string
	store.UploadFunc = func(ctx context.Context, path, contentType string, content io.Reader) error {
		gotContentType = contentType
		return nil
	}
	documentRepo.CreateFunc = func(ctx context.Context, d *models.Document) error { return nil }

	_, err := svc.UploadDocument(context.Background(), models.UploadDocumentInput{
		FamilyID: "family-1", UploadedBy: "member-1", FileName: "file.bin", SizeBytes: 5, Content: strings.NewReader("hello"),
	})

	require.NoError(t, err)
	assert.Equal(t, "application/octet-stream", gotContentType)
}

func TestUploadDocumentValidatesParentFamily(t *testing.T) {
	svc, _, parentRepo, _ := newTestDocumentService(t)

	parentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Parent, error) {
		return &models.Parent{ID: id, FamilyID: "other-family"}, nil
	}

	parentID := "parent-1"
	_, err := svc.UploadDocument(context.Background(), models.UploadDocumentInput{
		FamilyID: "family-1", UploadedBy: "member-1", FileName: "file.txt", SizeBytes: 5, ParentID: &parentID,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUploadDocumentCleansUpStorageOnMetadataFailure(t *testing.T) {
	svc, documentRepo, _, store := newTestDocumentService(t)

	store.UploadFunc = func(ctx context.Context, path, contentType string, content io.Reader) error {
		return nil
	}
	deletedPath := ""
	store.DeleteFunc = func(ctx context.Context, path string) error {
		deletedPath = path
		return nil
	}
	documentRepo.CreateFunc = func(ctx context.Context, d *models.Document) error {
		return errors.New("db exploded")
	}

	_, err := svc.UploadDocument(context.Background(), models.UploadDocumentInput{
		FamilyID: "family-1", UploadedBy: "member-1", FileName: "file.txt", SizeBytes: 5, Content: strings.NewReader("hello"),
	})

	require.Error(t, err)
	assert.NotEmpty(t, deletedPath, "should clean up the orphaned storage object")
}

func TestDownloadDocument(t *testing.T) {
	svc, documentRepo, _, store := newTestDocumentService(t)

	documentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Document, error) {
		return &models.Document{ID: id, StoragePath: "family-1/abc-file.txt"}, nil
	}
	store.DownloadFunc = func(ctx context.Context, path string) (io.ReadCloser, error) {
		assert.Equal(t, "family-1/abc-file.txt", path)
		return io.NopCloser(strings.NewReader("contents")), nil
	}

	doc, content, err := svc.DownloadDocument(context.Background(), "doc-1")

	require.NoError(t, err)
	assert.Equal(t, "doc-1", doc.ID)
	body, err := io.ReadAll(content)
	require.NoError(t, err)
	assert.Equal(t, "contents", string(body))
}

func TestDeleteDocumentCleansUpStorageAfterMetadataDeleted(t *testing.T) {
	svc, documentRepo, _, store := newTestDocumentService(t)

	documentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Document, error) {
		return &models.Document{ID: id, StoragePath: "family-1/abc-file.txt"}, nil
	}
	deleteOrder := []string{}
	documentRepo.DeleteFunc = func(ctx context.Context, id string) error {
		deleteOrder = append(deleteOrder, "db")
		return nil
	}
	store.DeleteFunc = func(ctx context.Context, path string) error {
		deleteOrder = append(deleteOrder, "storage")
		return nil
	}

	err := svc.DeleteDocument(context.Background(), "doc-1")

	require.NoError(t, err)
	assert.Equal(t, []string{"db", "storage"}, deleteOrder, "DB row must be gone before storage cleanup, never the reverse")
}

func TestDeleteDocumentNotFound(t *testing.T) {
	svc, documentRepo, _, _ := newTestDocumentService(t)

	documentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Document, error) {
		return nil, repository.ErrDocumentNotFound
	}

	err := svc.DeleteDocument(context.Background(), "doc-1")

	assert.ErrorIs(t, err, repository.ErrDocumentNotFound)
}
