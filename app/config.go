package app

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

// AppConfig represents the application configuration
type AppConfig struct {
	AppEnv  string `env:"APP_ENV" env-default:"dev"`
	Host    string `env:"HOST" env-default:"localhost"`
	Port    int    `env:"PORT" env-default:"3000"`
	Metrics MetricsConfig

	// UseHttpin enables httpin integration for request parsing
	UseHttpin bool `env:"USE_HTTPIN" env-default:"false"`
}

// MetricsConfig represents metrics server configuration
type MetricsConfig struct {
	Enabled bool   `env:"METRICS_ENABLED" env-default:"false"`
	Mode    string `env:"METRICS_MODE" env-default:"separate"` // "combined" or "separate"
	Path    string `env:"METRICS_PATH" env-default:"/metrics"` // endpoint path for combined mode
	Host    string `env:"METRICS_HOST" env-default:"localhost"`
	Port    int    `env:"METRICS_PORT" env-default:"9090"`
}

// DefaultAppConfig returns the default application configuration,
// loading values from environment variables
func DefaultAppConfig() AppConfig {
	var appConfig AppConfig
	err := cleanenv.ReadEnv(&appConfig)
	if err != nil {
		slog.Error("Failed reading environment variables", "err", err)
	}
	return appConfig
}

// Validate validates the application configuration
func (c *AppConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", c.Port)
	}

	if c.Metrics.Enabled {
		// Validate metrics mode
		if c.Metrics.Mode != "combined" && c.Metrics.Mode != "separate" {
			return fmt.Errorf("invalid metrics mode: %s (must be 'combined' or 'separate')", c.Metrics.Mode)
		}

		// Only validate separate port if in separate mode
		if c.Metrics.Mode == "separate" {
			if c.Metrics.Port < 1 || c.Metrics.Port > 65535 {
				return fmt.Errorf("invalid metrics port: %d (must be between 1 and 65535)", c.Metrics.Port)
			}
			if c.Metrics.Port == c.Port {
				return fmt.Errorf("metrics port cannot be the same as application port: %d", c.Port)
			}
		}

		// Validate path for combined mode
		if c.Metrics.Mode == "combined" {
			if c.Metrics.Path == "" {
				return fmt.Errorf("metrics path cannot be empty in combined mode")
			}
			if c.Metrics.Path[0] != '/' {
				return fmt.Errorf("metrics path must start with '/': %s", c.Metrics.Path)
			}
		}
	}

	return nil
}
