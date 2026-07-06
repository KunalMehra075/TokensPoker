package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	// Clear every variable Load reads so we observe the built-in defaults.
	for _, k := range []string{
		"APP_ENV", "PORT", "MONGODB_URI", "MONGODB_DB", "JWT_SECRET",
		"JWT_EXPIRY_HOURS", "CORS_ORIGINS", "RATE_LIMIT_RPS", "RATE_LIMIT_BURST",
	} {
		t.Setenv(k, "") // t.Setenv restores the prior value after the test.
	}

	cfg := Load()

	if cfg.Env != "development" {
		t.Errorf("Env = %q, want development", cfg.Env)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.MongoURI != "mongodb://localhost:27017" {
		t.Errorf("MongoURI = %q", cfg.MongoURI)
	}
	if cfg.MongoDB != "freetokenspoker" {
		t.Errorf("MongoDB = %q", cfg.MongoDB)
	}
	if cfg.JWTExpiry != 24*7*time.Hour {
		t.Errorf("JWTExpiry = %v, want %v", cfg.JWTExpiry, 24*7*time.Hour)
	}
	if cfg.RateLimitRPS != 20 || cfg.RateLimitBurst != 40 {
		t.Errorf("rate limits = %d/%d, want 20/40", cfg.RateLimitRPS, cfg.RateLimitBurst)
	}
	wantOrigins := []string{"http://localhost:5173", "http://localhost:4173"}
	if !equalStrings(cfg.CORSOrigins, wantOrigins) {
		t.Errorf("CORSOrigins = %v, want %v", cfg.CORSOrigins, wantOrigins)
	}
	if cfg.IsProduction() {
		t.Error("IsProduction() should be false in development")
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("JWT_EXPIRY_HOURS", "48")
	t.Setenv("RATE_LIMIT_RPS", "5")
	t.Setenv("CORS_ORIGINS", "https://freetokenspoker.com, https://www.freetokenspoker.com ")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want 9090", cfg.Port)
	}
	if cfg.JWTExpiry != 48*time.Hour {
		t.Errorf("JWTExpiry = %v, want 48h", cfg.JWTExpiry)
	}
	if cfg.RateLimitRPS != 5 {
		t.Errorf("RateLimitRPS = %d, want 5", cfg.RateLimitRPS)
	}
	if !cfg.IsProduction() {
		t.Error("IsProduction() should be true when APP_ENV=production")
	}
	want := []string{"https://freetokenspoker.com", "https://www.freetokenspoker.com"}
	if !equalStrings(cfg.CORSOrigins, want) {
		t.Errorf("CORSOrigins = %v, want %v (trimmed)", cfg.CORSOrigins, want)
	}
}

func TestGetEnvIntFallsBackOnInvalid(t *testing.T) {
	t.Setenv("RATE_LIMIT_RPS", "not-a-number")
	cfg := Load()
	if cfg.RateLimitRPS != 20 {
		t.Errorf("invalid int should fall back to default 20, got %d", cfg.RateLimitRPS)
	}
}

func TestSplitAndTrimDropsEmpties(t *testing.T) {
	cases := map[string][]string{
		"":                {},
		"a":               {"a"},
		" a , b ,c":       {"a", "b", "c"},
		"a,,b,":           {"a", "b"},
		"  ,  ,  ":        {},
		"one, two,three ": {"one", "two", "three"},
	}
	for in, want := range cases {
		got := splitAndTrim(in)
		if !equalStrings(got, want) {
			t.Errorf("splitAndTrim(%q) = %v, want %v", in, got, want)
		}
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
