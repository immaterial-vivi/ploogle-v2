package main

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/urfave/negroni"
)

func BasicGuard(username string, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryUsername, queryPassword, ok := r.BasicAuth()
		if !ok || queryUsername != username || queryPassword != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="protected"`)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			log.Println("rejected unauthorized request to", r.URL.Path, "by", r.UserAgent())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireApiKey(apiKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqKey := r.Header.Get("Authorization")

		if reqKey != apiKey {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			log.Println("rejected request without api key to", r.URL.Path, "by", r.UserAgent())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequestLog(dbpool *pgxpool.Pool, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path, r.UserAgent())
		lrw := negroni.NewResponseWriter(w)

		next.ServeHTTP(lrw, r)
		statusCode := lrw.Status()

		requestInfo := RequestInfo{
			Status:    statusCode,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
		}

		LogRequest(dbpool, requestInfo)

		log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))
	}
}
