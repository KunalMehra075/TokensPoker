package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime configuration, loaded from environment variables.
// Sensible defaults let the service boot locally with zero setup.
type Config struct {
	Env            string
	Port           string
	MongoURI       string
	MongoDB        string
	JWTSecret      string
	JWTExpiry      time.Duration
	CORSOrigins    []string
	RateLimitRPS   int
	RateLimitBurst int
}

// Load reads configuration from the environment, applying defaults.
func Load() *Config {
	cfg := &Config{
		Env:            getEnv("APP_ENV", "development"),
		Port:           getEnv("PORT", "8080"),
		MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:        getEnv("MONGODB_DB", "freetokenspoker"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-insecure-secret-change-me"),
		JWTExpiry:      time.Duration(getEnvInt("JWT_EXPIRY_HOURS", 24*7)) * time.Hour,
		CORSOrigins:    splitAndTrim(getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:4173")),
		RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 20),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 40),
	}
	return cfg
}

// IsProduction reports whether the service is running in production mode.
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
