package main

import (
	"context"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	Log *zap.Logger
}

type Proccessor struct {
	Log *zap.Logger
}

func main() {

	log, _ := zap.NewDevelopment()

	router := chi.NewRouter()

	handler := Handler{
		Log: log,
	}

	router.Post("/", handler.HandleEvent)

	log.Info("Started server on localhost:3000")
	err := http.ListenAndServe("localhost:3000", router)
	if err != nil {
		log.Fatal("Failed starting server", zap.Error(err))
	}

}

func (handler *Handler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	handler.Log.Info("Handling event")
	ce, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		handler.Log.Error("Failed creating event from HTTP Request.", zap.Error(err))
	}
	proccessor := Proccessor{
		Log: handler.Log,
	}
	proccessor.Receive(r.Context(), ce)

}

func (proccessor *Proccessor) Receive(ctx context.Context, event *cloudevents.Event) {
	proccessor.Log.Info("Proccessing event", zap.Any("event", event))
}
