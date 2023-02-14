package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handle struct {
	Log *zap.Logger
}

func (handle *Handle) Demo(w http.ResponseWriter, r *http.Request) {
	handle.Log.Debug("Demo Log")
	fmt.Fprintf(w, "Hello, Demo")
}

func Routes(r *chi.Mux, handle Handle) {
	r.Get("/demo", handle.Demo)
}
