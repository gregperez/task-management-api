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

func TestConfig_GetEnvNonExistent(t *testing.T) {
	defaultValue := "default"
	value := getEnv("NON_EXISTENT_ENV_VAR", defaultValue)
	if value != defaultValue {
		t.Errorf("Expected default value %s, got %s", defaultValue, value)
	}
}