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

	userRepo := userRepo.NewUserRepository(db)
	userUsecase := userUsecase.UserUsecase{
		Repository: userRepo,
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     configs.RedisAddr,
		Password: "",
		DB:       0,
	})

	sessionRepo := sessionRepo.NewSessionRedisRepo(redisClient)

	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRepo)

	sessionDelivery := sessionDelivery.NewSessionDelivery(sessionUsecase)

	userHandler := userDelivery.UserHandler{
		UserUsecase:    userUsecase,
		SessionHandler: sessionDelivery,
	}

	mux := mux.NewRouter().PathPrefix(configs.ApiUrl).Subrouter()
	mux.HandleFunc("/users", userHandler.Create).Methods("POST")

	siteHandler := middleware.CORS(mux)

	fmt.Println("start server at 9000")
	log.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
