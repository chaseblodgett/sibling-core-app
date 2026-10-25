package repository

import (
	"github.com/doug-martin/goqu/v9"
	// Registers the "postgres" dialect (e.g. $1, $2, ... placeholders
	// instead of goqu's generic "?"). Required as a side-effect import —
	// goqu.Dialect("postgres") silently falls back to the generic dialect
	// without it.
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

// dialect is the shared goqu SQL dialect used to build queries across all
// repositories.
var dialect = goqu.Dialect("postgres")
