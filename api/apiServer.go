package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/friends/storage"

	"github.com/gorilla/mux"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CORS middleware")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			fmt.Println("options")
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func StartApiServer() {
	mux := mux.NewRouter()
	db := storage.NewMapDB()
	userService := UserService{
		db: db,
	}
	mux.HandleFunc("/", userService.login)
	mux.HandleFunc("/reg", userService.reginster)
	mux.HandleFunc("/cookie", userService.testCookie)

	siteHandler := CORS(mux)

	log.Fatal(http.ListenAndServe(":9000", siteHandler))
}
