package config

import (
	"os"
	"time"
)

// Config holds the application configuration.
type Config struct {
	Port      string
	JWTSecret string
	JWTExpiry time.Duration
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "api-quest-secret-key-2026"),
		JWTExpiry: 1 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
