package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	profileDelivery "github.com/friends/internal/pkg/profile/delivery"
	profileRepo "github.com/friends/internal/pkg/profile/repository"
	profileUsecase "github.com/friends/internal/pkg/profile/usecase"
	sessionDelivery "github.com/friends/internal/pkg/session/delivery"
	sessionRepo "github.com/friends/internal/pkg/session/repository"
	sessionUsecase "github.com/friends/internal/pkg/session/usecase"
	userDelivery "github.com/friends/internal/pkg/user/delivery"
	userRepo "github.com/friends/internal/pkg/user/repository"
	userUsecase "github.com/friends/internal/pkg/user/usecase"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	logrus "github.com/sirupsen/logrus"
)

func StartApiServer() {
	db, err := sql.Open(configs.Postgres, configs.DataSourceNamePostgres)
	if err != nil {
		logrus.Error(fmt.Errorf("postgres not available: %w", err))
		fmt.Println("db doesnt work", err)
		return
	}
	err = db.Ping()
	if err != nil {
		logrus.Error(fmt.Errorf("no connection with db: %w", err))
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
		logrus.Error(fmt.Errorf("Session repostiory doen't work: %w", err))
		return
	}

	profRepo := profileRepo.NewProfileRepository(db)
	profUsecase := profileUsecase.NewProfileUsecase(profRepo)

	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRepo)

	userHandler := userDelivery.NewUserHandler(userUsecase, sessionUsecase, profUsecase)

	sessionDelivery := sessionDelivery.NewSessionDelivery(sessionUsecase, userUsecase)

	profDelivery := profileDelivery.NewProfileDelivery(profUsecase, sessionUsecase)

	authChecker := middleware.NewAuthChecker(sessionUsecase)

	mux := mux.NewRouter().PathPrefix(configs.ApiUrl).Subrouter()
	mux.HandleFunc("/users", userHandler.Create).Methods("POST")
	mux.HandleFunc("/users", userHandler.Delete).Methods("DELETE")
	mux.HandleFunc("/sessions", sessionDelivery.Create).Methods("POST")
	mux.HandleFunc("/sessions", sessionDelivery.Delete).Methods("DELETE")
	mux.Handle("/profiles", authChecker.Check(profDelivery.Get)).Methods("GET")
	mux.Handle("/profiles", authChecker.Check(profDelivery.Update)).Methods("PUT")

	accessLogHandler := middleware.AccessLog(mux)
	corsHandler := middleware.CORS(accessLogHandler)
	siteHandler := middleware.Panic(corsHandler)

	logrus.Info("starting server at port ", configs.Port)
	logrus.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
