package config

import (
	"time"

	"github.com/spf13/viper"
)

func setDefaults() {
	// Server
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 120*time.Second)
	viper.SetDefault("server.dev_mode", false)

	// Database
	viper.SetDefault("database.path", "./zoomgate.db")

	// Auth
	viper.SetDefault("auth.jwt_secret", "change-me-in-production")
	viper.SetDefault("auth.token_expiry", 24*time.Hour)
	viper.SetDefault("auth.admin_username", "admin")
	viper.SetDefault("auth.admin_password", "admin123")
	viper.SetDefault("auth.encryption_key", "change-me-32-byte-encryption-key")

	// Routing
	viper.SetDefault("routing.default_strategy", "priority")
	viper.SetDefault("routing.max_retries", 2)
	viper.SetDefault("routing.health_check_interval", 30*time.Second)

	// Security - Prompt Detection
	viper.SetDefault("security.prompt_detection.enabled", true)
	viper.SetDefault("security.prompt_detection.mode", "warn")
	viper.SetDefault("security.prompt_detection.threshold", 0.7)
	viper.SetDefault("security.prompt_detection.layers", []string{"regex", "heuristic", "structural"})

	// Security - Rate Limit
	viper.SetDefault("security.rate_limit.default_rpm", 60)
	viper.SetDefault("security.rate_limit.default_tpm", 100000)

	// Security - Data Masking
	viper.SetDefault("security.data_masking.enabled", true)
	viper.SetDefault("security.data_masking.patterns", []string{"email", "phone", "credit_card"})

	// Logging
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "console")
}
