package middleware

import (
	"log"
	"msng/internal/service"
	"net/http"
	"strings"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Starting: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Println("Done", time.Since(start))
	})
}

func RequestSizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(auth service.AuthTokenVerifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h := r.Header.Get("Authorization")
		if h == "" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "missing autherization header", http.StatusForbidden)
		}

		header_parts := strings.SplitN(h, " ", 2)
		if len(header_parts) != 2 || !strings.EqualFold(header_parts[0], "Bearer") || header_parts[1] == "" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "invalid authenticate header", http.StatusForbidden)
		}

		_, err := auth.VerifyAccessToken(header_parts[1])
		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "invalid access token", http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}
