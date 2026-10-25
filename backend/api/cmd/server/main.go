package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"sibling-core/api/internal/di"
	"sibling-core/api/internal/handler"
)

func main() {
	// Loads backend/api/.env when present; in production, real env vars are
	// already set and no .env file exists, so a missing file is not an error.
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found, using existing environment variables")
	}

	ctx := context.Background()

	container, err := di.Build(ctx)
	if err != nil {
		log.Fatalf("failed to build dependency container: %v", err)
	}
	defer container.Close()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS("*"))

	handler.RegisterRoutes(e, container.Handlers)

	addr := ":" + container.Config.Port
	log.Printf("starting server on %s (env=%s)", addr, container.Config.Env)
	if err := e.Start(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
