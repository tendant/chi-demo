package app

import (
	"log/slog"

	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/mikejav/gosts"
	"github.com/tendant/cors"

	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
)

// App represents the application with its router, configuration, and middleware stack
type App struct {
	R               *chi.Mux
	Config          AppConfig
	Logger          *slog.Logger
	HTTPLogger      *httplog.Logger
	middlewareStack *MiddlewareStack

	// Internal configuration (not directly exposed)
	corsOptions     *cors.Options
	metricsRecorder metricsMiddleware.Middleware
	hstsConfig      *gosts.Info
}

// DefaultCorsOptions returns CORS options with sensible defaults
func DefaultCorsOptions() *cors.Options {
	return &cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}
}

// NewApp creates a new application with the given options.
// By default, it creates a chi router with the default middleware stack.
//
// Example:
//
//	app := NewApp(
//	    WithPort(8080),
//	    WithDefaultCORS(),
//	    WithMetrics(true),
//	)
func NewApp(opts ...Option) *App {
	app := &App{
		R:               chi.NewRouter(),
		Config:          DefaultAppConfig(),
		middlewareStack: DefaultMiddlewareStack().Build(),
	}

	// Apply user options
	for _, opt := range opts {
		opt(app)
	}

	// Configure httpin if enabled
	if app.Config.UseHttpin {
		httpin_integration.UseGochiURLParam("path", chi.URLParam)
	}

	// Apply the middleware stack to the router
	if app.middlewareStack != nil {
		app.middlewareStack.Apply(app.R)
	}

	// Log version information
	slog.Info("Application initialized", "commit", Commit, "timestamp", Timestamp)

	return app
}

// DefaultApp creates an application with recommended defaults:
//   - Default configuration from environment
//   - Default logger based on APP_ENV
//   - Default HTTP logger
//   - CORS enabled with permissive defaults
//   - Metrics disabled (enable with WithMetrics or WithMetricsSeparate)
//   - HSTS enabled
//   - Httpin integration enabled
//
// This is the recommended way to create an app for most use cases.
func DefaultApp() *App {
	return NewApp(
		WithConfig(DefaultAppConfig()),
		WithLogger(DefaultLogger()),
		WithHTTPLogger(DefaultHTTPLogger()),
		WithDefaultCORS(),
		WithDefaultHSTS(),
		WithHttpin(true),
	)
}

// Run starts the HTTP server and blocks until shutdown.
// It handles graceful shutdown on SIGINT/SIGTERM signals.
func (app *App) Run() error {
	srv := &Server{App: app}
	return srv.Run()
}
