package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/tendant/chi-demo/app"
)

// Custom middleware example
func customAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple auth example - in production, use real authentication
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "Missing API key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("=== Middleware Stack Demo ===\n")

	// Example 1: Using default stack
	fmt.Println("Example 1: Default middleware stack")
	stack1 := app.DefaultMiddlewareStack()
	fmt.Println("Middleware in stack:")
	for _, name := range stack1.List() {
		fmt.Println("  ", name)
	}
	fmt.Println()

	// Example 2: Customizing the stack
	fmt.Println("Example 2: Custom middleware stack with auth")
	stack2 := app.DefaultMiddlewareStack().
		InsertAfter("request-id", "auth", customAuthMiddleware).
		Remove("hsts").
		Disable("metrics")

	fmt.Println("Modified middleware stack:")
	for _, name := range stack2.List() {
		fmt.Println("  ", name)
	}
	fmt.Println()

	// Example 3: Building a minimal stack
	fmt.Println("Example 3: Minimal stack")
	stack3 := app.MinimalMiddlewareStack().
		Append("no-cache", middleware.NoCache)

	fmt.Println("Minimal stack:")
	for _, name := range stack3.List() {
		fmt.Println("  ", name)
	}
	fmt.Println()

	// Create app with custom stack
	fmt.Println("Starting server with custom middleware stack...")
	myApp := app.NewApp(
		app.WithPort(3000),
		app.WithMiddlewareStack(stack2.Build()),
		app.WithDefaultCORS(),
	)

	// Register routes
	app.RegisterDefaultRoutes(myApp.R)
	app.RegisterVersionRoutes(myApp.R)
	app.RegisterHealthzRoutes(myApp.R)

	myApp.R.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello! Your middleware stack is working!\n"))
	})

	fmt.Println("Server running on http://localhost:3000")
	fmt.Println("Try: curl -H 'X-API-Key: test' http://localhost:3000/hello")
	if err := myApp.Run(); err != nil {
		log.Fatal(err)
	}
}
