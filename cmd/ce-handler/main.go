package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Host string `env:"HOST" env-default:"localhost"`
	Port int    `env:"PORT" env-default:"3000"`
}

type Handle struct {
	Log *zap.Logger
}

type Proccess struct {
	Log *zap.Logger
}

func main() {

	var cfg Config
	cleanenv.ReadEnv(&cfg)

	log, _ := zap.NewDevelopment()

	router := chi.NewRouter()

	handle := Handle{
		Log: log,
	}
	router.Post("/", handle.HandleEvent)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &http.Server{Addr: addr, Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed starting server", zap.Error(err))
		}
	}()

	log.Info(fmt.Sprintf("Started server on %s...", addr))

	// Capturing signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Waiting for SIGINT (kill -2)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Failed shutdown server", zap.Error(err))
	}
	log.Info("Server exited")

}

func (handler *Handle) HandleEvent(w http.ResponseWriter, r *http.Request) {
	handler.Log.Info("Handling event")
	ce, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		handler.Log.Error("Failed creating event from HTTP Request.", zap.Error(err))
	}
	proccess := Proccess{
		Log: handler.Log,
	}
	proccess.Receive(r.Context(), ce)

}

func (proccess *Proccess) Receive(ctx context.Context, event *cloudevents.Event) {
	proccess.Log.Info("Proccessing event", zap.Any("event", event))
}
