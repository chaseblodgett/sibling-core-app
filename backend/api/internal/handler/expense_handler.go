package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// ExpenseHandler exposes HTTP endpoints for expense tracking, splitting,
// and balance calculation.
type ExpenseHandler struct {
	expenseService service.ExpenseService
}

// NewExpenseHandler constructs an ExpenseHandler.
func NewExpenseHandler(expenseService service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

// RegisterRoutes wires expense endpoints onto the given group.
func (h *ExpenseHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/families/:familyId/expenses", h.CreateExpense)
	g.GET("/families/:familyId/expenses", h.ListExpenses)
	g.GET("/families/:familyId/balances", h.GetBalances)
	g.GET("/expenses/:id", h.GetExpense)
	g.DELETE("/expenses/:id", h.DeleteExpense)
	g.POST("/expenses/splits/:splitId/settle", h.SettleSplit)
}

// CreateExpense logs a new expense and its split among members.
func (h *ExpenseHandler) CreateExpense(c *echo.Context) error {
	var req models.CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// ListExpenses returns all expenses for a family.
func (h *ExpenseHandler) ListExpenses(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// GetExpense returns a single expense by id.
func (h *ExpenseHandler) GetExpense(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// DeleteExpense removes an expense.
func (h *ExpenseHandler) DeleteExpense(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// SettleSplit marks a member's share of an expense as settled.
func (h *ExpenseHandler) SettleSplit(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// GetBalances returns each member's running net balance for a family.
func (h *ExpenseHandler) GetBalances(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}
