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
	"github.com/go-chi/httplog"
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
	R      *chi.Mux
	Config AppConfig
	Slog   *slog.Logger
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
		JSON: false,
	})

	httpin.UseGochiURLParam("path", chi.URLParam)

	r := chi.NewRouter()

	app := &App{
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
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// AllowedOrigins: []string{"http://localhost:3000", "https://*.example.com"},
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
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

	return app
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

func RoutesDefault(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
}

func Default() *App {

	app := DefaultWithoutRoutes()
	RoutesDefault(app.R)
	RoutesVersion(app.R)
	RoutesHealthz(app.R)

	return app
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

	// Serve metrics.
	metricsAddr := fmt.Sprintf("%s:%d", app.Config.Metrics.Host, app.Config.Metrics.Port)
	metricsServer := &http.Server{Addr: metricsAddr, Handler: promhttp.Handler()}
	go func() {
		slog.Info("metrics listening at", "Addr", metricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed starting metrics server", "err", err)
		}
	}()

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
	if err := metricsServer.Shutdown(ctx); err != nil {
		slog.Error("Failed shutdown metrics server", "err", err)
	}
	slog.Info("Metrics Server exited")

}
