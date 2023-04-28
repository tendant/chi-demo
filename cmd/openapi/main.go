package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/ilyakaznacheev/cleanenv"
)

//go:embed static
var staticContent embed.FS

//go:embed api
var apiContent embed.FS

type Config struct {
	Port string `yaml:"port" env:"PORT" env-default:"3000"`
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
}

var cfg Config

func main() {
	cleanenv.ReadEnv(&cfg)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(middleware.BasicAuth("API Docs", map[string]string{"user": "Pu9yoo5Eemae3oekoh1e"}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "welcome")
	})

	fsys := fs.FS(staticContent)
	static, err := fs.Sub(fsys, "static")
	if err != nil {
		fmt.Println("Error add static/ folder", err)
		return
	}
	staticHandler := http.FileServer(http.FS(static))
	r.Handle("/docs/*", http.StripPrefix("/docs/", staticHandler))
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusSeeOther)
	})

	fapi := fs.FS(apiContent)
	api, err := fs.Sub(fapi, "api")
	if err != nil {
		fmt.Println("Error add api folder", err)
		return
	}
	apiHandler := http.FileServer(http.FS(api))
	r.Get("/api/", http.NotFound)
	r.Handle("/api/*", http.StripPrefix("/api/", apiHandler))

	server := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	http.ListenAndServe(server, r)
}
