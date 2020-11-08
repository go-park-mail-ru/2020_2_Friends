package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/friends/configs"
	cartDelivery "github.com/friends/internal/pkg/cart/delivery"
	cartRepo "github.com/friends/internal/pkg/cart/repository"
	cartUsecase "github.com/friends/internal/pkg/cart/usecase"
	"github.com/friends/internal/pkg/middleware"
	partnerDelivery "github.com/friends/internal/pkg/partner/delivery"
	profileDelivery "github.com/friends/internal/pkg/profile/delivery"
	profileRepo "github.com/friends/internal/pkg/profile/repository"
	profileUsecase "github.com/friends/internal/pkg/profile/usecase"
	sessionDelivery "github.com/friends/internal/pkg/session/delivery"
	sessionRepo "github.com/friends/internal/pkg/session/repository"
	sessionUsecase "github.com/friends/internal/pkg/session/usecase"
	userDelivery "github.com/friends/internal/pkg/user/delivery"
	userRepo "github.com/friends/internal/pkg/user/repository"
	userUsecase "github.com/friends/internal/pkg/user/usecase"
	vendorDelivery "github.com/friends/internal/pkg/vendors/delivery"
	vendorRepo "github.com/friends/internal/pkg/vendors/repository"
	vendorUsecase "github.com/friends/internal/pkg/vendors/usecase"
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

	vendRepo := vendorRepo.NewVendorRepository(db)
	vendUsecase := vendorUsecase.NewVendorUsecase(vendRepo)

	userHandler := userDelivery.NewUserHandler(userUsecase, sessionUsecase, profUsecase)

	sessionDelivery := sessionDelivery.NewSessionDelivery(sessionUsecase, userUsecase)

	profDelivery := profileDelivery.NewProfileDelivery(profUsecase, sessionUsecase)

	vendDelivery := vendorDelivery.NewVendorDelivery(vendUsecase)

	cartRepo := cartRepo.NewCartRepository(db)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepo, vendRepo)
	cartDelivery := cartDelivery.NewCartDelivery(cartUsecase)

	partnerDelivery := partnerDelivery.New(userUsecase, sessionUsecase, vendUsecase)

	authChecker := middleware.NewAuthChecker(sessionUsecase)
	accessRighsChecker := middleware.NewAccessRightsChecker(userUsecase)

	mux := mux.NewRouter().PathPrefix(configs.ApiUrl).Subrouter()
	mux.HandleFunc("/users", userHandler.Create).Methods("POST")
	mux.HandleFunc("/users", userHandler.Delete).Methods("DELETE")
	mux.HandleFunc("/sessions", sessionDelivery.Create).Methods("POST")
	mux.HandleFunc("/sessions", sessionDelivery.Delete).Methods("DELETE")
	mux.Handle("/profiles", authChecker.Check(profDelivery.Get)).Methods("GET")
	mux.Handle("/profiles", authChecker.Check(profDelivery.Update)).Methods("PUT")
	mux.Handle("/profiles/avatars", authChecker.Check(profDelivery.UpdateAvatar)).Methods("PUT")
	mux.HandleFunc("/vendors", vendDelivery.GetAll).Methods("GET")
	mux.HandleFunc("/vendors/{id}", vendDelivery.GetVendor).Methods("GET")
	mux.Handle("/vendors", accessRighsChecker.AccessRightsCheck(partnerDelivery.CreateVendor, configs.AdminRole)).Methods("POST")
	mux.Handle("/vendors/{id}", accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendor, configs.AdminRole)).Methods("PUT")
	mux.Handle("/vendors/{id}/pictures", accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendorPicture, configs.AdminRole)).Methods("PUT")
	mux.Handle("/vendors/products", accessRighsChecker.AccessRightsCheck(partnerDelivery.AddProductToVendor, configs.AdminRole)).Methods("POST")
	mux.Handle("/vendors/products/{id}", accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductOnVendor, configs.AdminRole)).Methods("PUT")
	mux.Handle("/vendors/products/{id}", accessRighsChecker.AccessRightsCheck(partnerDelivery.DeleteProductFromVendor, configs.AdminRole)).Methods("DELETE")
	mux.Handle("/vendors/products/{id}/avatars", accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductPicture, configs.AdminRole)).Methods("PUT")
	mux.Handle("/carts", authChecker.Check(cartDelivery.AddToCart)).Methods("PUT")
	mux.Handle("/carts", authChecker.Check(cartDelivery.RemoveFromCart)).Methods("DELETE")
	mux.Handle("/carts", authChecker.Check(cartDelivery.GetCart)).Methods("GET")
	mux.HandleFunc("/partners", partnerDelivery.Create).Methods("POST")

	accessLogHandler := middleware.AccessLog(mux)
	corsHandler := middleware.CORS(accessLogHandler)
	siteHandler := middleware.Panic(corsHandler)

	logrus.Info("starting server at port ", configs.Port)
	logrus.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
