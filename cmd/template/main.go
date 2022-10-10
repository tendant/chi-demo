package main

import (
	//...

	"crypto/sha1"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	"golang.org/x/oauth2"
)

func main() {
	r := chi.NewRouter()

	// AUTH
	// define options
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			log.Printf("secret reader: %s", id)
			return "secret", nil
		}),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "http://localhost:8080/realms/demo",
		URL:            "http://localhost:3000",
		AvatarStore:    avatar.NewLocalFS("/tmp"),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			log.Printf("Claims: %s", claims)
			// allow only dev_* names
			return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
		}),
	}
	// create auth service with providers
	service := auth.NewService(options)

	c := auth.Client{
		// Cid:     os.Getenv("OIDC_CLIENT_ID"),
		// Csecret: os.Getenv("OIDC_CLIENT_SECRET"),
		Cid:     "golang-demo",
		Csecret: "NPwPxyBPRg8Otu4Gh32TpdLsrAJ9SGRq",
	}

	service.AddCustomProvider("oidc", c, provider.CustomHandlerOpt{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8080/realms/demo/protocol/openid-connect/auth",
			TokenURL: "http://localhost:8080/realms/demo/protocol/openid-connect/token",
		},
		InfoURL: "http://localhost:8080/realms/demo/protocol/openid-connect/userinfo",
		MapUserFn: func(data provider.UserData, _ []byte) token.User {
			userInfo := token.User{
				ID: "localhost_" + token.HashID(sha1.New(),
					data.Value("username")),
				Name: data.Value("nickname"),
			}
			return userInfo
		},
		Scopes: []string{"account"},
	})
	// retrieve auth middleware
	m := service.Middleware()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.With(m.Auth).Get("/private", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("*PRIVATE*"))
	}) // protected api
	// setup auth routes
	authRoutes, avaRoutes := service.Handlers()
	r.Mount("/auth", authRoutes)  // add auth handlers
	r.Mount("/avatar", avaRoutes) // add avatar handler

	log.Print("Started server on localhost:3000...")
	http.ListenAndServe("localhost:3000", r)
}
