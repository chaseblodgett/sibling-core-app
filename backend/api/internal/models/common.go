package models

// ErrorResponse is the standard JSON shape returned for handler errors.
type ErrorResponse struct {
	Message string `json:"message"`
}
