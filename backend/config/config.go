package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	AllowedOrigins string
	Environment   string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Port:          GetEnv("PORT", "8080"),
		DatabaseURL:   GetEnv("DATABASE_URL", "postgres://user:password@postgres:5432/vconfdb?sslmode=disable"),
		JWTSecret:     GetEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		AllowedOrigins: GetEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		Environment:   GetEnv("ENVIRONMENT", "development"),
	}
}

// GetEnv retrieves an environment variable or returns a default value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvAsInt retrieves an environment variable as an integer or returns a default value
func GetEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// GetEnvAsBool retrieves an environment variable as a boolean or returns a default value
func GetEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}