package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv          string
	Port            string
	DatabaseURL     string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	JWTSecret       string
	JWTExpiresHours int
	AdminEmail      string
	AdminPassword   string
	AdminName       string
	AdminPhone      string
	CORSOrigins     []string
}

// Load reads configuration from environment variables with default fallbacks.
func Load() Config {
	expiresHours, _ := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))

	return Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "booking"),
		DBPassword:      getEnv("DB_PASSWORD", "booking"),
		DBName:          getEnv("DB_NAME", "special_booking"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me"),
		JWTExpiresHours: expiresHours,
		AdminEmail:      getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword:   getEnv("ADMIN_PASSWORD", "Admin@1234"),
		AdminName:       getEnv("ADMIN_NAME", "System Admin"),
		AdminPhone:      getEnv("ADMIN_PHONE", "0000000000"),
		CORSOrigins:     splitCSV(getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173")),
	}
}

// getEnv retrieves an environment variable or returns a default fallback value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// splitCSV splits a comma-separated string into a slice of trimmed strings, ignoring empty items.
func splitCSV(s string) []string {
	var results []string
	for _, item := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			results = append(results, trimmed)
		}
	}
	return results
}
