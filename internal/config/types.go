package config

import "time"

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Security SecurityConfig
	Features FeatureConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
}

type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	URL string
}

type SecurityConfig struct {
	JWTSecret        string
	ApprovalRequired bool
}

type FeatureConfig struct {
	AI         bool
	Monitoring bool
	Billing    bool
	Network    bool
}
