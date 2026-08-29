package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.App.Name == "" {
		t.Fatal("expected application name")
	}
	if cfg.Server.Port == "" {
		t.Fatal("expected server port")
	}
}

func TestValidateProductionRequiresDatabase(t *testing.T) {
	cfg := Load()
	cfg.App.Environment = "production"
	cfg.Database.URL = ""
	cfg.Security.JWTSecret = "a-secure-test-secret-with-more-than-32-chars"
	if err := Validate(cfg); err == nil {
		t.Fatal("expected database validation error")
	}
}

func TestValidateRejectsWeakProductionSecret(t *testing.T) {
	cfg := Load()
	cfg.App.Environment = "production"
	cfg.Database.URL = "postgres://localhost/ftn"
	cfg.Security.JWTSecret = "short"
	if err := Validate(cfg); err == nil {
		t.Fatal("expected authentication secret validation error")
	}
}
