// Package di wires together config, the database pool, repositories,
// services, and handlers using plain constructor-based dependency injection.
// There's no reflection-based DI library here on purpose: at this size,
// explicit constructors in one place are easier to trace than a container
// that resolves dependencies for you.
package di

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/config"
	"sibling-core/api/internal/handler"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
)

// Container holds every wired dependency needed to run the API server.
type Container struct {
	Config   *config.Config
	DB       *pgxpool.Pool
	Handlers *handler.Handlers
}

// Build loads config, connects to the database, and constructs the full
// repository -> service -> handler dependency graph.
func Build(ctx context.Context) (*Container, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	userRepo := repository.NewUserRepository(db)
	familyRepo := repository.NewFamilyRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	documentRepo := repository.NewDocumentRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	familyService := service.NewFamilyService(familyRepo, memberRepo, userRepo)
	taskService := service.NewTaskService(taskRepo)
	expenseService := service.NewExpenseService(expenseRepo)
	documentService := service.NewDocumentService(documentRepo)

	// Handlers
	handlers := &handler.Handlers{
		Auth:     handler.NewAuthHandler(authService),
		Family:   handler.NewFamilyHandler(familyService, authService),
		Task:     handler.NewTaskHandler(taskService),
		Expense:  handler.NewExpenseHandler(expenseService),
		Document: handler.NewDocumentHandler(documentService),
	}

	return &Container{
		Config:   cfg,
		DB:       db,
		Handlers: handlers,
	}, nil
}

// Close releases resources held by the container (e.g. the DB pool).
func (c *Container) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
}
