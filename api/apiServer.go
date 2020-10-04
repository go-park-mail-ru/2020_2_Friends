package api

import (
	"log"
	"net/http"

	"github.com/friends/storage"

	"github.com/gorilla/mux"
)

func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func StartApiServer() {
	r := mux.NewRouter()
	db := storage.NewMapDB()
	userService := UserService{
		db: db,
	}
	r.HandleFunc("/", userService.login)
	r.HandleFunc("/reg", userService.reginster)

	log.Fatal(http.ListenAndServe(":9000", r))
}
