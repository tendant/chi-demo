package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tendant/chi-demo/app"
	"golang.org/x/exp/slog"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsStd "github.com/slok/go-http-metrics/middleware/std"
)

type Server struct {
	Config Config
	R      *chi.Mux
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.Config.Server.Host, s.Config.Server.Port)
	server := &http.Server{Addr: addr, Handler: s.R}

	slog.Info("Started server.", "addr", addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed starting server", "err", err)
			os.Exit(-1)
		}
	}()

	// Serve metrics.
	metricsAddr := fmt.Sprintf("%s:%d", s.Config.Metrics.Host, s.Config.Metrics.Port)
	metricsServer := &http.Server{Addr: metricsAddr, Handler: promhttp.Handler()}
	go func() {
		slog.Info("metrics listening at", "Addr", metricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed starting metrics server", "err", err)
			// os.Exit(-1)
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

func Default(config Config) *Server {
	httpin.UseGochiURLParam("path", chi.URLParam)
	r := chi.NewRouter()

	s := &Server{
		R:      r,
		Config: config,
	}

	mdlw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	r.Use(metricsStd.HandlerProvider("", mdlw))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON: false,
		QuietDownRoutes: []string{
			"/healthz",
			"/healthz/ready",
		},
	})
	r.Use(httplog.RequestLogger(logger))

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

	return s
}

func Routes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, app.Commit)
	})

	r.Get("/version/timestamp", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, app.Timestamp)
	})
}
