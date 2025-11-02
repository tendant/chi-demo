package app

import (
	"log/slog"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/lmittmann/tint"
)

// LogConfig represents logger configuration
type LogConfig struct {
	Level      slog.Level
	Format     string // "text", "json", "tint"
	AddSource  bool
	TimeFormat string
}

// NewLogger creates a new slog logger with the given configuration
func NewLogger(config LogConfig) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	switch config.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "tint":
		handler = tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  config.AddSource,
			Level:      config.Level,
			TimeFormat: config.TimeFormat,
		})
	case "text":
		fallthrough
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// DefaultLogger returns a logger based on the APP_ENV environment variable.
// For development (APP_ENV=dev), it returns a colorful tint logger.
// For production, it returns a JSON logger.
func DefaultLogger() *slog.Logger {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	if appEnv == "dev" {
		return NewLogger(LogConfig{
			Level:      slog.LevelDebug,
			Format:     "tint",
			AddSource:  true,
			TimeFormat: time.DateTime,
		})
	}

	return NewLogger(LogConfig{
		Level:     slog.LevelInfo,
		Format:    "json",
		AddSource: false,
	})
}

// DefaultHTTPLogger returns an HTTP request logger with sensible defaults
func DefaultHTTPLogger() *httplog.Logger {
	return httplog.NewLogger("httplog", httplog.Options{
		JSON:             false,
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		QuietDownRoutes: []string{
			"/ping",
			"/healthz",
			"/healthz/ready",
		},
		QuietDownPeriod: 600 * time.Second,
	})
}
