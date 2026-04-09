package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	ExternalAPIKey string
	LocationSyncDays int
}

// Load reads configuration from environment variables.
// It attempts to load a .env file first, but does not fail if it's missing (for production).
func Load() *Config {
	_ = godotenv.Load() // Best practice: ignore error, env vars may come from Docker/K8s

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8088"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5433"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "location_demo"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		ExternalAPIKey: getEnv("EXTERNAL_API_KEY", ""),
		LocationSyncDays: getEnvAsInt("LOCATION_SYNC_DAYS", 30),
	}
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		var val int
		fmt.Sscanf(v, "%d", &val)
		return val
	}
	return fallback
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
