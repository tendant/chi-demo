package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/tendant/chi-demo/app"
	"github.com/tendant/chi-demo/middleware"
	"golang.org/x/exp/slog"
)

func main() {
	apiApp := app.Default()

	apiKeyConfig := middleware.ApiKeyConfig{
		APIKeys: map[string]string{
			"key1": "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", // echo -n "abc" |sha256sum
		},
	}
	apiKeyMiddleware, err := middleware.ApiKeyMiddleware(apiKeyConfig)
	if err != nil {
		slog.Error("Failed initialize API Key middleware", "err", err)
	} else {
		apiApp.R.Group(func(r chi.Router) {
			r.Use(apiKeyMiddleware)
			r.Route("/api", func(r chi.Router) {
				r.Get("/", handleApi)
			})
		})
		apiApp.Run()
	}
}

// curl -i -H "Authorization: abc" localhost:4000/api

func handleApi(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "API OK!")
}
