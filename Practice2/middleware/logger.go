package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			"request",
		)
		next.ServeHTTP(w, r)
	})
}