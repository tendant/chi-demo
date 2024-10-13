package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// Deprecated: Use middleware in `app` package
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing logging middleware")
		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		log.Printf("%q", dump)
		next.ServeHTTP(w, r)
		log.Print("Executing logging middleware again")
	})
}
