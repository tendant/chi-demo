package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server manages the HTTP server lifecycle including graceful shutdown
type Server struct {
	App           *App
	HTTPServer    *http.Server
	MetricsServer *http.Server

	// Shutdown timeout (default: 5 seconds)
	ShutdownTimeout time.Duration
}

// Run starts the HTTP server and handles graceful shutdown.
// It blocks until the server is shut down via signal (SIGINT, SIGTERM).
func (s *Server) Run() error {
	if s.ShutdownTimeout == 0 {
		s.ShutdownTimeout = 5 * time.Second
	}

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", s.App.Config.Host, s.App.Config.Port)
	s.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: s.App.R,
	}

	// Handle metrics based on mode
	if s.App.Config.Metrics.Enabled {
		switch s.App.Config.Metrics.Mode {
		case "combined":
			// Register metrics endpoint on main router
			s.App.R.Handle(s.App.Config.Metrics.Path, promhttp.Handler())
			slog.Info("Metrics enabled", "mode", "combined", "path", s.App.Config.Metrics.Path, "addr", addr)

		case "separate":
			// Start separate metrics server with custom path support
			metricsAddr := fmt.Sprintf("%s:%d", s.App.Config.Metrics.Host, s.App.Config.Metrics.Port)

			// Create a router for the metrics server to support custom paths
			metricsRouter := http.NewServeMux()
			metricsRouter.Handle(s.App.Config.Metrics.Path, promhttp.Handler())

			s.MetricsServer = &http.Server{
				Addr:    metricsAddr,
				Handler: metricsRouter,
			}

			go func() {
				slog.Info("Starting metrics server", "mode", "separate", "addr", metricsAddr, "path", s.App.Config.Metrics.Path)
				if err := s.MetricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("Metrics server error", "err", err)
				}
			}()

		default:
			slog.Warn("Invalid metrics mode", "mode", s.App.Config.Metrics.Mode)
		}
	}

	// Start main HTTP server
	slog.Info("Starting HTTP server", "addr", addr)
	go func() {
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "err", err)
		}
	}()

	// Wait for interrupt signal
	return s.waitForShutdown()
}

// waitForShutdown blocks until receiving SIGINT or SIGTERM, then gracefully shuts down
func (s *Server) waitForShutdown() error {
	// Create channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	sig := <-stop
	slog.Info("Received shutdown signal", "signal", sig)

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		slog.Error("HTTP server shutdown error", "err", err)
		return err
	}
	slog.Info("HTTP server stopped gracefully")

	// Shutdown metrics server if running
	if s.MetricsServer != nil {
		if err := s.MetricsServer.Shutdown(ctx); err != nil {
			slog.Error("Metrics server shutdown error", "err", err)
			return err
		}
		slog.Info("Metrics server stopped gracefully")
	}

	return nil
}

// Shutdown gracefully shuts down the server with the given context
func (s *Server) Shutdown(ctx context.Context) error {
	if s.HTTPServer != nil {
		if err := s.HTTPServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("http server shutdown: %w", err)
		}
	}

	if s.MetricsServer != nil {
		if err := s.MetricsServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("metrics server shutdown: %w", err)
		}
	}

	return nil
}
