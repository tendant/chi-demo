package handler

import (
	"fmt"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Handle struct {
	Log *zap.Logger
}

func (handle *Handle) Demo(w http.ResponseWriter, r *http.Request) {
	handle.Log.Debug("Demo Log")
	fmt.Fprintf(w, "Hello, Demo")
}

type QueryInput struct {
	Q string `in:"query=q"`
}

func (handle *Handle) Query(w http.ResponseWriter, r *http.Request) {
	handle.Log.Debug("Query...")
	q := r.Context().Value(httpin.Input).(*QueryInput)
	handle.Log.Debug("Q:", zap.Any("q", q))
	render.PlainText(w, r, "Query")
}

type DemoPostInput struct {
	Name string `in:"form=name"`
}

func (handle *Handle) DemoPost(w http.ResponseWriter, r *http.Request) {
	handle.Log.Debug("Form...")
	q := r.Context().Value(httpin.Input).(*DemoPostInput)
	handle.Log.Debug("Q:", zap.Any("q", q))
	render.PlainText(w, r, "DemoPost")
}

func Routes(r *chi.Mux, handle Handle) {
	r.Get("/demo", handle.Demo)
	r.With(httpin.NewInput(QueryInput{})).Get("/query", handle.Query)
	r.With(httpin.NewInput(DemoPostInput{})).Post("/post", handle.DemoPost)
}
