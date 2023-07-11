package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	zerologl "github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsStd "github.com/slok/go-http-metrics/middleware/std"
)

type App struct {
	R      *chi.Mux
	Config AppConfig
	Log    *zap.Logger
	Zlog   zerolog.Logger
	Slog   *slog.Logger
}

func Default() *App {

	// Configuration
	var appConfig AppConfig
	cleanenv.ReadEnv(&appConfig)

	// Logger
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	slogger := slog.New(jsonHandler)
	// textHandler := slog.NewTextHandler(os.Stdout, nil)
	// slogger := slog.New(textHandler)

	log, err := zap.NewDevelopment()
	if err != nil {
		panic("Can't initialize Zap log!")
	}

	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON: false,
	})

	zlog := zerologl.With().Str("service", "app")

	httpin.UseGochiURLParam("path", chi.URLParam)

	r := chi.NewRouter()

	app := &App{
		R:      r,
		Config: appConfig,
		Log:    log,
		Zlog:   zlog.Logger(),
		Slog:   slogger,
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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, Commit)
	})

	r.Get("/version/timestamp", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, Timestamp)
	})

	return app
}

func (app *App) Run() {
	addr := fmt.Sprintf("%s:%d", app.Config.Host, app.Config.Port)
	server := &http.Server{Addr: addr, Handler: app.R}

	app.Slog.Info("Started server.", "addr", addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Slog.Error("Failed starting server", "err", err)
		}
	}()

	// Serve metrics.
	metricsAddr := fmt.Sprintf("%s:%d", app.Config.Metrics.Host, app.Config.Metrics.Port)
	metricsServer := &http.Server{Addr: metricsAddr, Handler: promhttp.Handler()}
	go func() {
		app.Slog.Info("metrics listening at", "Addr", metricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Slog.Error("Failed starting metrics server", "err", err)
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
		app.Slog.Error("Failed shutdown server", "err", err)
	}
	app.Slog.Info("Server exited")
	if err := metricsServer.Shutdown(ctx); err != nil {
		app.Slog.Error("Failed shutdown metrics server", "err", err)
	}
	app.Slog.Info("Metrics Server exited")

}
