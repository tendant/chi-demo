package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lmittmann/tint"
	"github.com/mikejav/gosts"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsStd "github.com/slok/go-http-metrics/middleware/std"
)

type App struct {
	R             *chi.Mux
	Config        AppConfig
	Slog          *slog.Logger
	HttpLogger    *httplog.Logger
	CorsOptions   *cors.Options
	UseHttpin     bool
	EnableMetrics bool
}

type Option func(*App)

func DefaultAppConfig() AppConfig {
	var appConfig AppConfig
	err := cleanenv.ReadEnv(&appConfig)
	if err != nil {
		slog.Error("Failed reading environment variables", "err", err)
	}
	return appConfig
}

func DefaultHttpLogger() *httplog.Logger {
	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON:             false,
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		// Tags: map[string]string{
		// 	"version": "v1.0-81aa4244d9fc8076a",
		// 	"env":     "dev",
		// },
		QuietDownRoutes: []string{
			"/ping",
			"/healthz",
			"/healthz/ready",
		},
		QuietDownPeriod: 600 * time.Second,
		// SourceFieldName: "source",
	})
	return logger
}

type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func CustomCorsOptions(config CorsConfig) *cors.Options {
	return &cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   config.AllowedMethods,
		AllowedHeaders:   config.AllowedHeaders,
		ExposedHeaders:   config.ExposedHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	}
}

func DefaultCorsOptions() *cors.Options {
	return CustomCorsOptions(CorsConfig{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

func NewApp(opts ...Option) *App {
	server := &App{}
	for _, opt := range opts {
		opt(server)
	}
	if server.R == nil {
		server.R = chi.NewRouter()
	}
	if server.UseHttpin {
		httpin.UseGochiURLParam("path", chi.URLParam)
	}
	if server.EnableMetrics {
		mdlw := metricsMiddleware.New(metricsMiddleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{}),
		})

		server.R.Use(metricsStd.HandlerProvider("", mdlw))

	}
	server.R.Use(middleware.RequestID)
	server.R.Use(middleware.RealIP)
	server.R.Use(middleware.Recoverer)
	slog.Info("appVersion", "commit", Commit, "timestamp", Timestamp)
	server.R.Use(Version(Commit))

	if server.HttpLogger != nil {
		server.R.Use(httplog.RequestLogger(server.HttpLogger))
	}

	if server.CorsOptions != nil {
		server.R.Use(cors.Handler(*server.CorsOptions))
	}

	// For Security
	server.R.Use(middleware.NoCache)

	// For Security: HSTS
	// config for hsts middleware
	hstsConf := &gosts.Info{
		MaxAge:               60 * 60 * 24,
		Expires:              time.Now().Add(24 * time.Hour),
		IncludeSubDomains:    true,
		SendPreloadDirective: false,
	}
	// middleware
	gosts.Configure(hstsConf)
	server.R.Use(gosts.Header)

	return server
}

func WithHttpin(useHttpin bool) Option {
	return func(s *App) {
		s.UseHttpin = useHttpin
	}
}

func WithReqLogger(logger *httplog.Logger) Option {
	return func(s *App) {
		s.HttpLogger = logger
	}
}

func WithAppConfig(config AppConfig) Option {
	return func(s *App) {
		s.Config = config
	}
}

func WithCors(corsOptions *cors.Options) Option {
	return func(s *App) {
		s.CorsOptions = corsOptions
	}
}

func WithRouter(router *chi.Mux) Option {
	return func(s *App) {
		s.R = router
	}
}

func WithMetrics(enable bool) Option {
	return func(s *App) {
		s.EnableMetrics = enable
	}
}

func DefaultWithoutRoutes() *App {

	// Configuration
	var appConfig AppConfig
	cleanenv.ReadEnv(&appConfig)

	// Logger
	appEnv := appConfig.AppEnv
	slog.Info("appEnv", "appEnv", appEnv)
	var slogger *slog.Logger
	if appEnv == "dev" {
		// create a new logger
		slogger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		}))

		// textHandler := slog.NewTextHandler(os.Stdout, nil)
		// slogger = slog.New(textHandler)
	} else {
		jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
		slogger = slog.New(jsonHandler)
	}
	slog.SetDefault(slogger)

	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON:             false,
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		// Tags: map[string]string{
		// 	"version": "v1.0-81aa4244d9fc8076a",
		// 	"env":     "dev",
		// },
		QuietDownRoutes: []string{
			"/ping",
			"/healthz",
			"/healthz/ready",
		},
		QuietDownPeriod: 600 * time.Second,
		// SourceFieldName: "source",
	})

	httpin.UseGochiURLParam("path", chi.URLParam)

	r := chi.NewRouter()

	server := &App{
		R:      r,
		Config: appConfig,
	}

	mdlw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	r.Use(metricsStd.HandlerProvider("", mdlw))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(httplog.RequestLogger(logger))
	// r.Use(middleware.Logger)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	server.CorsOptions = &cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// AllowedOrigins: []string{"http://localhost:3000", "https://*.example.com"},
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}
	r.Use(cors.Handler(*server.CorsOptions))
	r.Use(middleware.NoCache)

	// config for hsts middleware
	hstsConf := &gosts.Info{
		MaxAge:               60 * 60 * 24,
		Expires:              time.Now().Add(24 * time.Hour),
		IncludeSubDomains:    true,
		SendPreloadDirective: false,
	}
	// middleware
	gosts.Configure(hstsConf)
	r.Use(gosts.Header)

	return server
}

func RoutesVersion(r *chi.Mux) {
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, Commit)
	})

	r.Get("/version/timestamp", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, Timestamp)
	})
}

func RoutesHealthz(r *chi.Mux) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
}

func RoutesHealthzReady(r *chi.Mux) {
	r.Get("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
}

func RoutesDefault(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
}

func DefaultApp() *App {
	server := NewApp(
		WithAppConfig(DefaultAppConfig()),
		WithMetrics(true),
		WithCors(DefaultCorsOptions()),
		WithHttpin(true),
		WithMetrics(true),
		WithReqLogger(DefaultHttpLogger()),
	)
	return server
}

func Default() *App {

	server := DefaultWithoutRoutes()
	RoutesDefault(server.R)
	RoutesVersion(server.R)
	RoutesHealthz(server.R)

	return server
}

func (app *App) Run() {
	addr := fmt.Sprintf("%s:%d", app.Config.Host, app.Config.Port)
	server := &http.Server{Addr: addr, Handler: app.R}

	slog.Info("Started server.", "addr", addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed starting server", "err", err)
		}
	}()

	var metricsServer *http.Server
	if app.EnableMetrics {
		// Serve metrics.
		metricsAddr := fmt.Sprintf("%s:%d", app.Config.Metrics.Host, app.Config.Metrics.Port)
		metricsServer = &http.Server{Addr: metricsAddr, Handler: promhttp.Handler()}
		go func() {
			slog.Info("metrics listening at", "Addr", metricsAddr)
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("Failed starting metrics server", "err", err)
			}
		}()
	}

	// Capturing signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Waiting for SIGINT (kill -2)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed shutdown server", "err", err)
	}
	slog.Info("Server exited")
	if app.EnableMetrics && metricsServer != nil {
		if err := metricsServer.Shutdown(ctx); err != nil {
			slog.Error("Failed shutdown metrics server", "err", err)
		}
		slog.Info("Metrics Server exited")
	}

}
