package main

import (
	"log"
	"net/http"
)

func BasicGuard(next http.Handler, user string, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		validUsername := user
		validPassword := password
		username, passowrd, ok := r.BasicAuth()

		if !ok || username != validUsername || passowrd != validPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="protected"`)
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			log.Println("rejected unauthorized request to", r.URL.Path, "by", r.UserAgent())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireApiKey(next http.Handler, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqKey := r.Header.Get("x-ploogle-api-key")

		if reqKey != apiKey {
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			log.Println("rejected request without api key to", r.URL.Path, "by", r.UserAgent())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}
