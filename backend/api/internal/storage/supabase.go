package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const bucket = "documents"

// supabaseStorage talks to Supabase Storage's object API directly:
// https://supabase.com/docs/guides/storage — POST/GET/DELETE
// {projectURL}/storage/v1/object/{bucket}/{path}, authenticated with the
// service role key (bypasses the bucket's RLS entirely, which is fine
// since only this backend ever calls it).
type supabaseStorage struct {
	baseURL    string // {SUPABASE_URL}/storage/v1/object/documents
	serviceKey string
	httpClient *http.Client
}

// NewSupabaseStorage constructs a Storage backed by Supabase Storage.
// supabaseURL is the project URL (e.g. http://127.0.0.1:54321 locally).
func NewSupabaseStorage(supabaseURL, serviceRoleKey string) Storage {
	return &supabaseStorage{
		baseURL:    strings.TrimRight(supabaseURL, "/") + "/storage/v1/object/" + bucket,
		serviceKey: serviceRoleKey,
		httpClient: &http.Client{},
	}
}

func (s *supabaseStorage) objectURL(path string) string {
	return s.baseURL + "/" + strings.TrimLeft(path, "/")
}

func (s *supabaseStorage) authHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("apikey", s.serviceKey)
}

func (s *supabaseStorage) Upload(ctx context.Context, path, contentType string, content io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.objectURL(path), content)
	if err != nil {
		return fmt.Errorf("building upload request: %w", err)
	}
	s.authHeaders(req)
	req.Header.Set("Content-Type", contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("uploading to storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage upload returned %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (s *supabaseStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.objectURL(path), nil)
	if err != nil {
		return nil, fmt.Errorf("building download request: %w", err)
	}
	s.authHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading from storage: %w", err)
	}

	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrObjectNotFound
		}
		return nil, fmt.Errorf("storage download returned %d: %s", resp.StatusCode, string(body))
	}
	return resp.Body, nil
}

func (s *supabaseStorage) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.objectURL(path), nil)
	if err != nil {
		return fmt.Errorf("building delete request: %w", err)
	}
	s.authHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("deleting from storage: %w", err)
	}
	defer resp.Body.Close()

	// A missing object isn't an error for a delete — the end state (the
	// object is gone) is already true.
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage delete returned %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
