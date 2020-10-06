package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/friends/storage"

	"github.com/gorilla/mux"
)

const apiUrl = "/api/v1"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CORS middleware")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			fmt.Println("options")
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func StartApiServer() {
	mux := mux.NewRouter().PathPrefix(apiUrl).Subrouter()
	db := storage.NewUserMapDB()
	userService := UserService{
		db: db,
	}

	vendorService := VendorService{}

	sessionDB := storage.NewSessionMapDB()
	sessionService := SessionService{
		db:          sessionDB,
		userService: userService,
	}

	mux.HandleFunc("/reg", userService.reginster).Methods("POST")
	mux.HandleFunc("/cookie", userService.testCookie).Methods("GET")
	mux.HandleFunc("/vendors/{id}", vendorService.getVendor).Methods("GET")
	mux.HandleFunc("/session", sessionService.login).Methods("POST")
	mux.HandleFunc("/session", sessionService.logout).Methods("DELETE")

	siteHandler := CORS(mux)

	log.Fatal(http.ListenAndServe(":9000", siteHandler))
}
