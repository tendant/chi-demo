package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/tendant/chi-demo/app"
)

func main() {
	// Create app with default settings
	// This gives you: CORS, HSTS, logging, graceful shutdown
	myApp := app.DefaultApp()

	// Add your routes
	myApp.R.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "Welcome to chi-demo!")
	})

	myApp.R.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		slog.Info("Hello endpoint called", "name", name)
		render.JSON(w, r, map[string]string{
			"message": "Hello, " + name + "!",
		})
	})

	// Start the server (blocks until shutdown signal)
	slog.Info("Starting basic example on http://localhost:3000")
	myApp.Run()
}
