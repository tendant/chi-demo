package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("welcome2"))
		render.PlainText(w, r, "hello5")
	})
	http.ListenAndServe("localhost:3000", r)
}
