package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Routing  RoutingConfig  `mapstructure:"routing"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	DevMode      bool          `mapstructure:"dev_mode"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type AuthConfig struct {
	JWTSecret     string        `mapstructure:"jwt_secret"`
	TokenExpiry   time.Duration `mapstructure:"token_expiry"`
	AdminUsername string        `mapstructure:"admin_username"`
	AdminPassword string        `mapstructure:"admin_password"`
	EncryptionKey string        `mapstructure:"encryption_key"`
}

type RoutingConfig struct {
	DefaultStrategy     string        `mapstructure:"default_strategy"`
	MaxRetries          int           `mapstructure:"max_retries"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
}

type SecurityConfig struct {
	PromptDetection PromptDetectionConfig `mapstructure:"prompt_detection"`
	RateLimit       RateLimitConfig       `mapstructure:"rate_limit"`
	DataMasking     DataMaskingConfig     `mapstructure:"data_masking"`
}

type PromptDetectionConfig struct {
	Enabled     bool     `mapstructure:"enabled"`
	Mode        string   `mapstructure:"mode"` // block | warn | off
	Threshold   float64  `mapstructure:"threshold"`
	Layers      []string `mapstructure:"layers"`
	ExternalURL string   `mapstructure:"external_url"`
}

type RateLimitConfig struct {
	DefaultRPM int `mapstructure:"default_rpm"`
	DefaultTPM int `mapstructure:"default_tpm"`
}

type DataMaskingConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	Patterns []string `mapstructure:"patterns"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json | console
}

// Load loads the configuration from a file or environment variables
func Load(path string) (*Config, error) {
	// Set default values
	setDefaults()

	if path != "" {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
	}

	viper.SetEnvPrefix("ZOOMGATE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
		// Config file not found is OK, use defaults + env vars
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
