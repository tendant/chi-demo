package main

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
	"github.com/go-chi/render"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Host string `env:"HOST" env-default:"localhost"`
	Port int    `env:"PORT" env-default:"3000"`
}

func main() {

	// Configuration
	var cfg Config
	cleanenv.ReadEnv(&cfg)

	// Logger
	log, _ := zap.NewDevelopment()

	router := chi.NewRouter()

	handle := Handle{
		Log: log,
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
	router.With(httpin.NewInput(GetInput{})).Get("/get", handle.HandleGet)
	router.With(httpin.NewInput(PostInput{})).Post("/body", handle.HandlePost)
	router.With(httpin.NewInput(FormInput{})).Post("/form", handle.HandleForm)

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

type Handle struct {
	Log *zap.Logger
}

type GetInput struct {
	Id          string `in:"query=id"`
	Name        string `in:"query=name"`
	notExported string `in:"query=notExported"`
}

func (handle *Handle) HandleGet(w http.ResponseWriter, r *http.Request) {
	handle.Log.Info("Handling get")
	query := r.Context().Value(httpin.Input).(*GetInput)
	handle.Log.Info("query", zap.String("id", query.Id))
	handle.Log.Info(fmt.Sprintf("query: %+v", query))
}

type FormInput struct {
	Name        string `in:"form=name"`
	Descritpion string `in:"form=description"`
}

func (handle *Handle) HandleForm(w http.ResponseWriter, r *http.Request) {
	handle.Log.Info("Handling form")
	form := r.Context().Value(httpin.Input).(*FormInput)
	handle.Log.Info(fmt.Sprintf("form: %+v", form))
}

type PostBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PostInput struct {
	Payload *PostBody `in:"body=json"`
}

func (handle *Handle) HandlePost(w http.ResponseWriter, r *http.Request) {
	handle.Log.Info("Handling request")
	process := Process{
		Log: handle.Log,
	}
	body := r.Context().Value(httpin.Input).(*PostInput)
	handle.Log.Info(fmt.Sprintf("body: %+v", body))
	process.Receive(body.Payload)

}

type Process struct {
	Log *zap.Logger
}

func (process *Process) Receive(payload *PostBody) {
	process.Log.Info(fmt.Sprintf("Processing payload: %+v", payload))
}
