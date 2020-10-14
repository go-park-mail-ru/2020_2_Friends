package middleware

import (
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
			return
		}
		next.ServeHTTP(w, r)
	})
}
