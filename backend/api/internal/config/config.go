// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"os"
)

// Config holds all runtime configuration for the API server.
type Config struct {
	Port         string
	Env          string
	DatabaseURL  string // Postgres connection string (Supabase)
	SupabaseURL  string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	JWTSecret    string
}

// Load reads configuration from environment variables, applying sane
// defaults for local development. Call this once at startup.
func Load() (*Config, error) {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("APP_ENV", "development"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:    os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
