package app

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/mikejav/gosts"
	"github.com/tendant/cors"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsStd "github.com/slok/go-http-metrics/middleware/std"
)

// Option is a functional option for configuring an App
type Option func(*App)

// WithConfig sets the application configuration
func WithConfig(config AppConfig) Option {
	return func(a *App) {
		a.Config = config
	}
}

// WithPort sets the HTTP server port
func WithPort(port int) Option {
	return func(a *App) {
		a.Config.Port = port
	}
}

// WithHost sets the HTTP server host
func WithHost(host string) Option {
	return func(a *App) {
		a.Config.Host = host
	}
}

// WithLogger sets a custom slog logger
func WithLogger(logger *slog.Logger) Option {
	return func(a *App) {
		a.Logger = logger
		// Set as default logger
		slog.SetDefault(logger)
	}
}

// WithHTTPLogger sets the HTTP request logger and enables it in the middleware stack
func WithHTTPLogger(logger *httplog.Logger) Option {
	return func(a *App) {
		a.HTTPLogger = logger

		// Enable http-logger in the middleware stack
		if a.middlewareStack != nil {
			for i, item := range a.middlewareStack.items {
				if item.Name == "http-logger" {
					a.middlewareStack.items[i].Enabled = true
					a.middlewareStack.items[i].Middleware = httplog.RequestLogger(logger)
					break
				}
			}
		}
	}
}

// WithLogLevel sets the log level for the default logger
func WithLogLevel(level slog.Level) Option {
	return func(a *App) {
		if a.Logger == nil {
			a.Logger = NewLogger(LogConfig{Level: level})
			slog.SetDefault(a.Logger)
		}
	}
}

// WithMiddlewareStack sets a custom middleware stack, completely overriding the default
func WithMiddlewareStack(stack *MiddlewareStack) Option {
	return func(a *App) {
		a.middlewareStack = stack
	}
}

// WithCORS enables and configures CORS middleware
func WithCORS(opts *cors.Options) Option {
	return func(a *App) {
		a.corsOptions = opts

		// Enable CORS in the middleware stack
		if a.middlewareStack != nil {
			for i, item := range a.middlewareStack.items {
				if item.Name == "cors" {
					a.middlewareStack.items[i].Enabled = true
					a.middlewareStack.items[i].Middleware = cors.Handler(*opts)
					break
				}
			}
		}
	}
}

// WithDefaultCORS enables CORS with default settings
func WithDefaultCORS() Option {
	return WithCORS(DefaultCorsOptions())
}

// WithMetrics enables and configures Prometheus metrics middleware in separate server mode.
// This maintains backward compatibility - metrics will run on a separate port (default: 9090).
// For combined mode (metrics on same server), use WithMetricsCombined() instead.
func WithMetrics(enabled bool) Option {
	return func(a *App) {
		if !enabled {
			return
		}

		// Create metrics recorder
		mdlw := metricsMiddleware.New(metricsMiddleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{}),
		})
		a.metricsRecorder = mdlw

		// Update config to enable metrics server in SEPARATE mode (backward compatible)
		a.Config.Metrics.Enabled = true
		a.Config.Metrics.Mode = "separate" // Keep old behavior

		// Enable metrics in the middleware stack
		if a.middlewareStack != nil {
			for i, item := range a.middlewareStack.items {
				if item.Name == "metrics" {
					a.middlewareStack.items[i].Enabled = true
					a.middlewareStack.items[i].Middleware = metricsStd.HandlerProvider("", mdlw)
					break
				}
			}
		}
	}
}

// WithMetricsCombined enables metrics on the main application server at the specified path.
// This is simpler than separate mode as it only requires one port.
// Default path is "/metrics" if not specified via WithMetricsPath().
func WithMetricsCombined() Option {
	return func(a *App) {
		// Create metrics recorder
		mdlw := metricsMiddleware.New(metricsMiddleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{}),
		})
		a.metricsRecorder = mdlw

		// Update config to enable metrics in COMBINED mode
		a.Config.Metrics.Enabled = true
		a.Config.Metrics.Mode = "combined"

		// Enable metrics in the middleware stack
		if a.middlewareStack != nil {
			for i, item := range a.middlewareStack.items {
				if item.Name == "metrics" {
					a.middlewareStack.items[i].Enabled = true
					a.middlewareStack.items[i].Middleware = metricsStd.HandlerProvider("", mdlw)
					break
				}
			}
		}
	}
}

// WithMetricsSeparatePort enables metrics on a separate dedicated server.
// This is useful for production environments where you want metrics isolated.
func WithMetricsSeparatePort(port int) Option {
	return func(a *App) {
		// Create metrics recorder
		mdlw := metricsMiddleware.New(metricsMiddleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{}),
		})
		a.metricsRecorder = mdlw

		// Update config to enable metrics in SEPARATE mode
		a.Config.Metrics.Enabled = true
		a.Config.Metrics.Mode = "separate"
		a.Config.Metrics.Port = port

		// Enable metrics in the middleware stack
		if a.middlewareStack != nil {
			for i, item := range a.middlewareStack.items {
				if item.Name == "metrics" {
					a.middlewareStack.items[i].Enabled = true
					a.middlewareStack.items[i].Middleware = metricsStd.HandlerProvider("", mdlw)
					break
				}
			}
		}
	}
}

// WithMetricsPath sets the endpoint path for metrics in combined mode.
// Default is "/metrics". Only applies when using combined mode.
func WithMetricsPath(path string) Option {
	return func(a *App) {
		a.Config.Metrics.Path = path
	}
}

// WithMetricsMode explicitly sets the metrics mode.
// Valid values: "combined" or "separate"
func WithMetricsMode(mode string) Option {
	return func(a *App) {
		a.Config.Metrics.Mode = mode
	}
}

// WithHSTS enables HTTP Strict Transport Security headers
func WithHSTS(config *gosts.Info) Option {
	return func(a *App) {
		a.hstsConfig = config
		gosts.Configure(config)

		// Enable HSTS in the middleware stack
		if a.middlewareStack != nil {
			for i, item := range a.middlewareStack.items {
				if item.Name == "hsts" {
					a.middlewareStack.items[i].Enabled = true
					a.middlewareStack.items[i].Middleware = gosts.Header
					break
				}
			}
		}
	}
}

// WithDefaultHSTS enables HSTS with default configuration
func WithDefaultHSTS() Option {
	return WithHSTS(&gosts.Info{
		MaxAge:               60 * 60 * 24, // 24 hours
		Expires:              time.Now().Add(24 * time.Hour),
		IncludeSubDomains:    true,
		SendPreloadDirective: false,
	})
}

// WithRouter sets a custom chi router
func WithRouter(router *chi.Mux) Option {
	return func(a *App) {
		a.R = router
	}
}

// WithHttpin enables httpin integration for request parsing
func WithHttpin(enabled bool) Option {
	return func(a *App) {
		a.Config.UseHttpin = enabled
	}
}
