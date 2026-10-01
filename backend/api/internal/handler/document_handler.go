package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// DocumentHandler exposes HTTP endpoints for uploading and listing family documents.
type DocumentHandler struct {
	documentService service.DocumentService
}

// NewDocumentHandler constructs a DocumentHandler.
func NewDocumentHandler(documentService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{documentService: documentService}
}

// RegisterRoutes wires document endpoints onto the given group.
func (h *DocumentHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/families/:familyId/documents", h.UploadDocument)
	g.GET("/families/:familyId/documents", h.ListDocuments)
	g.GET("/documents/:id", h.GetDocument)
	g.DELETE("/documents/:id", h.DeleteDocument)
}

// UploadDocument accepts a multipart file upload tied to a family.
func (h *DocumentHandler) UploadDocument(c *echo.Context) error {
	// Expects a multipart/form-data request with a "file" field.
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// ListDocuments returns all documents uploaded for a family.
func (h *DocumentHandler) ListDocuments(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// GetDocument returns metadata (and/or a signed URL) for a single document.
func (h *DocumentHandler) GetDocument(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// DeleteDocument removes a document.
func (h *DocumentHandler) DeleteDocument(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}
