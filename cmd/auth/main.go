package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
)

func main() {
	// define options
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			return "secret", nil
		}),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "my-test-app",
		URL:            "http://127.0.0.1:8080",
		AvatarStore:    avatar.NewLocalFS("/tmp"),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			// allow only dev_* names
			return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
		}),
	}

	// create auth service with providers
	service := auth.NewService(options)
	service.AddDirectProvider("local", provider.CredCheckerFunc(func(user, password string) (ok bool, err error) {
		ok, err = checkUserSomehow(user, password)
		return ok, err
	}))
	// service.AddProvider("github", "<Client ID>", "<Client Secret>")   // add github provider
	// service.AddProvider("facebook", "<Client ID>", "<Client Secret>") // add facebook provider

	// retrieve auth middleware
	m := service.Middleware()

	// setup http server
	router := chi.NewRouter()
	router.Get("/open", openRouteHandler)                      // open api
	router.With(m.Auth).Get("/private", protectedRouteHandler) // protected api

	// setup auth routes
	authRoutes, avaRoutes := service.Handlers()
	router.Mount("/auth", authRoutes)  // add auth handlers
	router.Mount("/avatar", avaRoutes) // add avatar handler

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

func checkUserSomehow(user string, password string) (bool, error) {
	log.Printf("Login with: %s, %s.", user, password)
	if user == "user1" && password == "pass1" {
		return true, nil
	}
	return false, errors.New("Incorrect username and/or password!")
}

func openRouteHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "Open Route")
}

func protectedRouteHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "*Protected Route*")
}
