package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// RegisterDefaultRoutes registers a simple root route that returns HTTP 200 OK
func RegisterDefaultRoutes(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
}

// RegisterVersionRoutes registers routes that return version information
func RegisterVersionRoutes(r chi.Router) {
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, Commit)
	})

	r.Get("/version/timestamp", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, Timestamp)
	})
}

// RegisterHealthzRoutes registers health check routes
func RegisterHealthzRoutes(r chi.Router) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})

	r.Get("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
}
