package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Validate(cfg Config) error {
	if strings.TrimSpace(cfg.App.Name) == "" {
		return errors.New("app name must not be empty")
	}
	if cfg.App.Environment != "development" && cfg.App.Environment != "staging" && cfg.App.Environment != "production" {
		return fmt.Errorf("unsupported FTN_ENV %q", cfg.App.Environment)
	}
	port, err := strconv.Atoi(cfg.Server.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid server port %q", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout <= 0 || cfg.Server.WriteTimeout <= 0 || cfg.Server.IdleTimeout <= 0 || cfg.Server.ShutdownTimeout <= 0 {
		return errors.New("server timeouts must be positive")
	}
	if cfg.Database.MaxOpenConns < 1 || cfg.Database.MaxIdleConns < 0 || cfg.Database.MaxIdleConns > cfg.Database.MaxOpenConns {
		return errors.New("invalid database pool settings")
	}
	if cfg.App.Environment == "production" {
		if strings.TrimSpace(cfg.Database.URL) == "" {
			return errors.New("database URL is required in production")
		}
		if len(cfg.Security.JWTSecret) < 32 {
			return errors.New("authentication secret must be at least 32 characters in production")
		}
	}
	return nil
}
