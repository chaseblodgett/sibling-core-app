package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
)

// ExpenseHandler exposes HTTP endpoints for expense tracking, splitting,
// and balance calculation. Every route requires the requester to belong
// to the relevant family; there's no further role restriction — any
// member can manage any expense, matching Task's "shared board" model.
type ExpenseHandler struct {
	expenseService service.ExpenseService
	familyService  service.FamilyService
	authService    service.AuthService
}

// NewExpenseHandler constructs an ExpenseHandler.
func NewExpenseHandler(expenseService service.ExpenseService, familyService service.FamilyService, authService service.AuthService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService, familyService: familyService, authService: authService}
}

// RegisterRoutes wires expense endpoints onto the given group, all behind
// middleware.RequireAuth.
func (h *ExpenseHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.POST("/families/:familyId/expenses", h.CreateExpense, auth)
	g.GET("/families/:familyId/expenses", h.ListExpenses, auth)
	g.GET("/families/:familyId/balances", h.GetBalances, auth)
	g.GET("/expenses/:id", h.GetExpense, auth)
	g.GET("/expenses/:id/splits", h.ListSplits, auth)
	g.DELETE("/expenses/:id", h.DeleteExpense, auth)
	g.POST("/expenses/splits/:splitId/settle", h.SettleSplit, auth)
}

// CreateExpense logs a new expense and its split among members, optionally
// tied to a parent (parent_id in the body).
func (h *ExpenseHandler) CreateExpense(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	requester, err := h.familyService.RequireMembership(ctx, user.ID, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	var req models.CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	paidBy := requester.ID
	if req.PaidBy != nil && *req.PaidBy != "" {
		paidBy = *req.PaidBy
	}

	splits := make([]models.SplitInput, len(req.Splits))
	for i, s := range req.Splits {
		splits[i] = models.SplitInput{MemberID: s.MemberID, ShareCents: s.ShareCents}
	}

	expense, err := h.expenseService.CreateExpense(ctx, models.CreateExpenseInput{
		FamilyID:    familyID,
		PaidBy:      paidBy,
		Description: req.Description,
		AmountCents: req.AmountCents,
		Currency:    req.Currency,
		Category:    req.Category,
		ParentID:    req.ParentID,
		Splits:      splits,
	})
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusCreated, expense)
}

// ListExpenses returns a family's expenses. An optional ?parent_id query
// param filters to expenses tied to that parent.
func (h *ExpenseHandler) ListExpenses(c *echo.Context) error {
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

	expenses, err := h.expenseService.ListExpenses(ctx, familyID, parentID)
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, expenses)
}

// GetBalances returns each member's running net balance for a family.
func (h *ExpenseHandler) GetBalances(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	balances, err := h.expenseService.GetBalances(ctx, familyID)
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, balances)
}

// GetExpense returns a single expense by id, if the requester belongs to its family.
func (h *ExpenseHandler) GetExpense(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	expense, err := h.expenseService.GetExpense(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, expense.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, expense)
}

// ListSplits returns an expense's per-member splits.
func (h *ExpenseHandler) ListSplits(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	expense, err := h.expenseService.GetExpense(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, expense.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	splits, err := h.expenseService.ListSplits(ctx, expense.ID)
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, splits)
}

// DeleteExpense removes an expense (and its splits, via cascade).
func (h *ExpenseHandler) DeleteExpense(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	expense, err := h.expenseService.GetExpense(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, expense.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	if err := h.expenseService.DeleteExpense(ctx, expense.ID); err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.NoContent(http.StatusNoContent)
}

// SettleSplit marks a member's share of an expense as settled. Membership
// is checked against the split's own expense/family before mutating.
func (h *ExpenseHandler) SettleSplit(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	splitID := c.Param("splitId")

	// We only have the split id, not its family, so resolve it via its
	// expense before authorizing — never mutate before that check passes.
	split, err := h.expenseService.GetSplit(ctx, splitID)
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	expense, err := h.expenseService.GetExpense(ctx, split.ExpenseID)
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, expense.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	updated, err := h.expenseService.SettleSplit(ctx, splitID)
	if err != nil {
		status, msg := mapExpenseError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, updated)
}

func mapExpenseError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, repository.ErrExpenseNotFound), errors.Is(err, repository.ErrSplitNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
