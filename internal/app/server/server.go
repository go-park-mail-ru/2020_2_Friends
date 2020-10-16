package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	sessionDelivery "github.com/friends/internal/pkg/session/delivery"
	sessionRepo "github.com/friends/internal/pkg/session/repository"
	sessionUsecase "github.com/friends/internal/pkg/session/usecase"
	userDelivery "github.com/friends/internal/pkg/user/delivery"
	userRepo "github.com/friends/internal/pkg/user/repository"
	userUsecase "github.com/friends/internal/pkg/user/usecase"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func StartApiServer() {
	db, err := sql.Open(configs.Postgres, configs.DataSourceNamePostgres)
	if err != nil {
		fmt.Println("db doesnt work", err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("db doesnt work", err)
		return
	}

	userRepo := userRepo.NewUserRepository(db)
	userUsecase := userUsecase.NewUserUsecase(userRepo)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     configs.RedisAddr,
		Password: "",
		DB:       0,
	})

	sessionRepo, err := sessionRepo.NewSessionRedisRepo(redisClient)
	if err != nil {
		fmt.Println("redis doesnt work")
		return
	}

	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRepo)

	userHandler := userDelivery.NewUserHandler(userUsecase, sessionUsecase)

	sessionDelivery := sessionDelivery.NewSessionDelivery(sessionUsecase, userUsecase)

	mux := mux.NewRouter().PathPrefix(configs.ApiUrl).Subrouter()
	mux.HandleFunc("/users", userHandler.Create).Methods("POST")
	mux.HandleFunc("/users", userHandler.Delete).Methods("DELETE")
	mux.HandleFunc("/sessions", sessionDelivery.Create).Methods("POST")
	mux.HandleFunc("/sessions", sessionDelivery.Delete).Methods("DELETE")

	corsHandler := middleware.CORS(mux)
	siteHandler := middleware.Panic(corsHandler)

	fmt.Println("start server at 9000")
	log.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
