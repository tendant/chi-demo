package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/render"
	"github.com/tendant/chi-demo/app"
)

func main() {
	// Create app with httpin enabled
	myApp := app.NewApp(
		app.WithPort(3000),
		app.WithDefaultCORS(),
		app.WithHttpin(true),
	)

	// Register routes
	myApp.R.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "Input Handling Examples\n\nEndpoints:\n- GET /query?name=John&age=30\n- GET /users?is_member=true&age_range=18&age_range=60\n- POST /form (application/x-www-form-urlencoded)\n- POST /json (application/json)")
	})

	// Query parameters
	myApp.R.With(httpin.NewInput(QueryInput{})).Get("/query", HandleQuery)

	// Query with arrays and booleans
	myApp.R.With(httpin.NewInput(ListUsersInput{})).Get("/users", HandleListUsers)

	// Form data
	myApp.R.With(httpin.NewInput(FormInput{})).Post("/form", HandleForm)

	// JSON body
	myApp.R.With(httpin.NewInput(JSONInput{})).Post("/json", HandleJSON)

	slog.Info("Input handling examples available at http://localhost:3000")
	myApp.Run()
}

// QueryInput demonstrates basic query parameter parsing
type QueryInput struct {
	Name string `in:"query=name"`
	Age  int    `in:"query=age"`
}

func HandleQuery(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*QueryInput)
	slog.Info("Query parameters", "name", input.Name, "age", input.Age)
	render.JSON(w, r, map[string]interface{}{
		"name": input.Name,
		"age":  input.Age,
	})
}

// ListUsersInput demonstrates arrays and boolean query parameters
type ListUsersInput struct {
	IsMember bool  `in:"query=is_member"`
	AgeRange []int `in:"query=age_range"`
}

func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)
	slog.Info("List users", "is_member", input.IsMember, "age_range", input.AgeRange)
	render.JSON(w, r, map[string]interface{}{
		"is_member": input.IsMember,
		"age_range": input.AgeRange,
	})
}

// FormInput demonstrates form data parsing
type FormInput struct {
	Name        string `in:"form=name"`
	Description string `in:"form=description"`
}

func HandleForm(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*FormInput)
	slog.Info("Form data", "name", input.Name, "description", input.Description)
	render.JSON(w, r, map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
	})
}

// JSONInput demonstrates JSON body parsing
type JSONInput struct {
	Payload *JSONBody `in:"body=json"`
}

type JSONBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []string `json:"tags"`
}

func HandleJSON(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*JSONInput)
	slog.Info("JSON body", "payload", fmt.Sprintf("%+v", input.Payload))
	render.JSON(w, r, input.Payload)
}
