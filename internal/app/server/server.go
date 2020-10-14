package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	userDelivery "github.com/friends/internal/pkg/user/delivery"
	userRepo "github.com/friends/internal/pkg/user/repository"
	userUsecase "github.com/friends/internal/pkg/user/usecase"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func StartApiServer() {
	db, err := sql.Open(configs.Postgres, configs.DataSourceNamePostgres)
	if err != nil {
		fmt.Println("db doesnt work", err)
		return
	}

	repo := userRepo.NewUserRepository(db)
	userUsecase := userUsecase.UserUsecase{
		Repository: repo,
	}

	userHandler := userDelivery.UserHandler{
		UserUsecase: userUsecase,
	}

	mux := mux.NewRouter().PathPrefix(configs.ApiUrl).Subrouter()
	mux.HandleFunc("/users", userHandler.Create).Methods("POST")

	siteHandler := middleware.CORS(mux)

	fmt.Println("start server at 9000")
	log.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
