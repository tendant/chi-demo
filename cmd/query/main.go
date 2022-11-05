package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

// Check document for more examples of directives
// https://ggicci.github.io/httpin/directives/query
// https://ggicci.github.io/httpin/directives/header
// https://ggicci.github.io/httpin/directives/form
// https://ggicci.github.io/httpin/directives/body
// https://ggicci.github.io/httpin/directives/required
// https://ggicci.github.io/httpin/directives/default
type ListUsersInput struct {
	IsMember bool  `in:"query=is_member"`
	AgeRange []int `in:"query=age_range[],age_range"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(httpin.Input).(*ListUsersInput)
	fmt.Printf("input: %#v\n", input)
}

func main() {
	router := chi.NewRouter()
	router.With(
		httpin.NewInput(ListUsersInput{}),
	).Get("/users", ListUsers)

	r, _ := http.NewRequest("GET", "/users?is_member=1&age_range=18&age_range=60", nil)

	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, r)
}
