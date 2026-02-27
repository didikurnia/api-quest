package config

import (
	"crypto/rand"
	"encoding/hex"
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
		JWTSecret: getEnv("JWT_SECRET", generateRandomSecret()),
		JWTExpiry: 1 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func generateRandomSecret() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
