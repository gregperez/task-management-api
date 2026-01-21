package config

import (
	"testing"
	"time"
)

func TestConfig_Load(t *testing.T) {
	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}
	if cfg.DatabaseURL == "" {
		t.Error("Expected non-empty DatabaseURL")
	}
	if cfg.JWTSecret == "" {
		t.Error("Expected non-empty JWTSecret")
	}
	if cfg.TokenExpiry != 2*time.Hour {
		t.Errorf("Expected default TokenExpiry 2 hours, got %v", cfg.TokenExpiry)
	}
}

func TestConfig_GetEnvTestVar(t *testing.T) {
	t.Setenv("TEST_VAR", "test_value")
	value := getEnv("TEST_VAR", "default_value")
	if value != "test_value" {
		t.Errorf("Expected value test_value, got %s", value)
	}
}