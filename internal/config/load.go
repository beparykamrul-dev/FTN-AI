package config

import "time"

func Load() Config {
	return Config{
		App: AppConfig{
			Name:        envString("FTN_APP_NAME", "FTN Enterprise"),
			Environment: envString("FTN_ENV", "production"),
			Version:     envString("FTN_VERSION", "dev"),
		},
		Server: ServerConfig{
			Host:            envString("FTN_HOST", "127.0.0.1"),
			Port:            envString("FTN_PORT", "8080"),
			ReadTimeout:     time.Duration(envInt("FTN_READ_TIMEOUT_SEC", 30)) * time.Second,
			WriteTimeout:    time.Duration(envInt("FTN_WRITE_TIMEOUT_SEC", 30)) * time.Second,
			IdleTimeout:     time.Duration(envInt("FTN_IDLE_TIMEOUT_SEC", 60)) * time.Second,
			ShutdownTimeout: time.Duration(envInt("FTN_SHUTDOWN_TIMEOUT_SEC", 15)) * time.Second,
		},
		Database: DatabaseConfig{
			URL:             envString("FTN_DATABASE_URL", ""),
			MaxOpenConns:    envInt("FTN_DB_MAX_OPEN", 25),
			MaxIdleConns:    envInt("FTN_DB_MAX_IDLE", 10),
			ConnMaxLifetime: time.Duration(envInt("FTN_DB_MAX_LIFETIME_SEC", 300)) * time.Second,
			ConnMaxIdleTime: time.Duration(envInt("FTN_DB_MAX_IDLE_TIME_SEC", 60)) * time.Second,
		},
		Redis: RedisConfig{URL: envString("FTN_REDIS_URL", "redis://127.0.0.1:6379")},
		Security: SecurityConfig{
			JWTSecret:        envString("FTN_JWT_SECRET", ""),
			ApprovalRequired: envBool("FTN_APPROVAL_REQUIRED", true),
		},
		Features: FeatureConfig{
			AI:         envBool("FTN_FEATURE_AI", true),
			Monitoring: envBool("FTN_FEATURE_MONITORING", true),
			Billing:    envBool("FTN_FEATURE_BILLING", true),
			Network:    envBool("FTN_FEATURE_NETWORK", true),
		},
	}
}
