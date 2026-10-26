package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	"sibling-core/api/internal/storage"
)

// DocumentHandler exposes HTTP endpoints for uploading, downloading, and
// listing family documents. Every route requires the requester to belong
// to the relevant family; there's no further role restriction, matching
// Task/Expense's shared-board model.
type DocumentHandler struct {
	documentService service.DocumentService
	familyService   service.FamilyService
	authService     service.AuthService
}

// NewDocumentHandler constructs a DocumentHandler.
func NewDocumentHandler(documentService service.DocumentService, familyService service.FamilyService, authService service.AuthService) *DocumentHandler {
	return &DocumentHandler{documentService: documentService, familyService: familyService, authService: authService}
}

// RegisterRoutes wires document endpoints onto the given group, all
// behind middleware.RequireAuth.
func (h *DocumentHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.POST("/families/:familyId/documents", h.UploadDocument, auth)
	g.GET("/families/:familyId/documents", h.ListDocuments, auth)
	g.GET("/documents/:id", h.GetDocument, auth)
	g.GET("/documents/:id/download", h.DownloadDocument, auth)
	g.DELETE("/documents/:id", h.DeleteDocument, auth)
}

// UploadDocument accepts a multipart/form-data request with a "file"
// field and an optional "parent_id" field, tied to a family.
func (h *DocumentHandler) UploadDocument(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	requester, err := h.familyService.RequireMembership(ctx, user.ID, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "missing \"file\" in multipart form"})
	}
	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "couldn't read uploaded file"})
	}
	defer file.Close()

	var parentID *string
	if v := c.FormValue("parent_id"); v != "" {
		parentID = &v
	}

	doc, err := h.documentService.UploadDocument(ctx, models.UploadDocumentInput{
		FamilyID:    familyID,
		UploadedBy:  requester.ID,
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		SizeBytes:   fileHeader.Size,
		ParentID:    parentID,
		Content:     file,
	})
	if err != nil {
		status, msg := mapDocumentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusCreated, doc)
}

// ListDocuments returns a family's documents. An optional ?parent_id
// query param filters to documents tied to that parent.
func (h *DocumentHandler) ListDocuments(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	var parentID *string
	if v := c.QueryParam("parent_id"); v != "" {
		parentID = &v
	}

	documents, err := h.documentService.ListDocuments(ctx, familyID, parentID)
	if err != nil {
		status, msg := mapDocumentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, documents)
}

// GetDocument returns a document's metadata, if the requester belongs to its family.
func (h *DocumentHandler) GetDocument(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	doc, err := h.documentService.GetDocument(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapDocumentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, doc.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, doc)
}

// DownloadDocument streams a document's file content.
func (h *DocumentHandler) DownloadDocument(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	// Check metadata + membership before ever touching storage — an
	// unauthorized request shouldn't cost a storage round-trip, and
	// definitely shouldn't fetch the bytes before being allowed to.
	doc, err := h.documentService.GetDocument(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapDocumentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, doc.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	_, content, err := h.documentService.DownloadDocument(ctx, doc.ID)
	if err != nil {
		status, msg := mapDocumentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	defer content.Close()

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, doc.FileName))
	return c.Stream(http.StatusOK, doc.ContentType, content)
}

// DeleteDocument removes a document (metadata row and storage object).
func (h *DocumentHandler) DeleteDocument(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	doc, err := h.documentService.GetDocument(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapDocumentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, doc.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	if err := h.documentService.DeleteDocument(ctx, doc.ID); err != nil {
		status, msg := mapDocumentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.NoContent(http.StatusNoContent)
}

func mapDocumentError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, repository.ErrDocumentNotFound), errors.Is(err, storage.ErrObjectNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
