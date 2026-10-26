package storage

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testServiceKey = "test-service-role-key"

func TestUploadSendsCorrectRequest(t *testing.T) {
	var gotMethod, gotPath, gotContentType, gotAuth, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotContentType = r.Header.Get("Content-Type")
		gotAuth = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewSupabaseStorage(srv.URL, testServiceKey)
	err := s.Upload(context.Background(), "family-1/abc-report.pdf", "application/pdf", strings.NewReader("file bytes"))

	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/storage/v1/object/documents/family-1/abc-report.pdf", gotPath)
	assert.Equal(t, "application/pdf", gotContentType)
	assert.Equal(t, "Bearer "+testServiceKey, gotAuth)
	assert.Equal(t, "file bytes", gotBody)
}

func TestUploadSurfacesServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("disk full"))
	}))
	defer srv.Close()

	s := NewSupabaseStorage(srv.URL, testServiceKey)
	err := s.Upload(context.Background(), "x", "text/plain", strings.NewReader("data"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "disk full")
}

func TestDownloadReturnsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/storage/v1/object/documents/family-1/report.pdf", r.URL.Path)
		_, _ = w.Write([]byte("the contents"))
	}))
	defer srv.Close()

	s := NewSupabaseStorage(srv.URL, testServiceKey)
	rc, err := s.Download(context.Background(), "family-1/report.pdf")
	require.NoError(t, err)
	defer rc.Close()

	body, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Equal(t, "the contents", string(body))
}

func TestDownloadNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	s := NewSupabaseStorage(srv.URL, testServiceKey)
	_, err := s.Download(context.Background(), "missing")

	assert.ErrorIs(t, err, ErrObjectNotFound)
}

func TestDeleteTreatsNotFoundAsSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	s := NewSupabaseStorage(srv.URL, testServiceKey)
	err := s.Delete(context.Background(), "already-gone")

	assert.NoError(t, err)
}

func TestDeleteSurfacesRealError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("nope"))
	}))
	defer srv.Close()

	s := NewSupabaseStorage(srv.URL, testServiceKey)
	err := s.Delete(context.Background(), "x")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "nope")
}
