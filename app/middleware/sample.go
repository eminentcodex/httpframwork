package middleware

import (
	"log"
	"net/http"
)

func Sample(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Sample Middleware")
		next.ServeHTTP(w,r)
	})
}

