package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/friends/configs"
	cartDelivery "github.com/friends/internal/pkg/cart/delivery"
	cartRepo "github.com/friends/internal/pkg/cart/repository"
	cartUsecase "github.com/friends/internal/pkg/cart/usecase"
	csrfDelivery "github.com/friends/internal/pkg/csrf/delivery"
	csrfRepo "github.com/friends/internal/pkg/csrf/repository"
	csrfUsecase "github.com/friends/internal/pkg/csrf/usecase"
	"github.com/friends/internal/pkg/middleware"
	orderDelivery "github.com/friends/internal/pkg/order/delivery"
	orderRepo "github.com/friends/internal/pkg/order/repository"
	orderUsecase "github.com/friends/internal/pkg/order/usecase"
	partnerDelivery "github.com/friends/internal/pkg/partner/delivery"
	profileDelivery "github.com/friends/internal/pkg/profile/delivery"
	profileRepo "github.com/friends/internal/pkg/profile/repository"
	profileUsecase "github.com/friends/internal/pkg/profile/usecase"
	reviewDelivery "github.com/friends/internal/pkg/review/delivery"
	reviewRepository "github.com/friends/internal/pkg/review/repository"
	reviewUsecase "github.com/friends/internal/pkg/review/usecase"
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

	profDelivery := profileDelivery.NewProfileDelivery(profUsecase)

	vendDelivery := vendorDelivery.NewVendorDelivery(vendUsecase)

	cartRepo := cartRepo.NewCartRepository(db)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepo, vendRepo)
	cartDelivery := cartDelivery.NewCartDelivery(cartUsecase)

	partnerDelivery := partnerDelivery.New(userUsecase, profUsecase, sessionUsecase, vendUsecase)

	orderRepo := orderRepo.New(db)
	orderUsecase := orderUsecase.New(orderRepo, vendRepo)
	orderDelivery := orderDelivery.New(orderUsecase, vendUsecase)

	reviewRepository := reviewRepository.New(db)
	reviewUsecase := reviewUsecase.New(reviewRepository, orderRepo)
	reviewDelivery := reviewDelivery.New(reviewUsecase, vendUsecase)

	accessRighsChecker := middleware.NewAccessRightsChecker(userUsecase)

	csrfRepository, err := csrfRepo.New(redisClient)
	if err != nil {
		logrus.Error(fmt.Errorf("CSRF repository not created: %w", err))
		return
	}
	csrfUsecase := csrfUsecase.New(csrfRepository)
	csrfDelivery := csrfDelivery.New(csrfUsecase)

	authChecker := middleware.NewAuthChecker(sessionUsecase)
	csrfChecker := middleware.NewCSRFChecker(authChecker, csrfUsecase)

	mux := mux.NewRouter().PathPrefix(configs.ApiUrl).Subrouter()
	mux.HandleFunc("/users", userHandler.Create).Methods("POST")
	mux.HandleFunc("/users", userHandler.Delete).Methods("DELETE")
	mux.HandleFunc("/sessions", sessionDelivery.Create).Methods("POST")
	mux.HandleFunc("/sessions", sessionDelivery.Delete).Methods("DELETE")
	mux.Handle("/profiles", csrfChecker.Check(profDelivery.Get)).Methods("GET")
	mux.Handle("/profiles", csrfChecker.Check(profDelivery.Update)).Methods("PUT")
	mux.Handle("/profiles/avatars", csrfChecker.Check(profDelivery.UpdateAvatar)).Methods("PUT")
	mux.Handle("/profiles/addresses", csrfChecker.Check(profDelivery.UpdateAddresses)).Methods("PUT")
	mux.HandleFunc("/vendors", vendDelivery.GetAll).Methods("GET")
	mux.HandleFunc("/vendors/{id}", vendDelivery.GetVendor).Methods("GET")
	mux.Handle("/vendors", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.CreateVendor, configs.AdminRole))).Methods("POST")
	mux.Handle("/vendors/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendor, configs.AdminRole))).Methods("PUT")
	mux.Handle("/vendors/{id}/pictures", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendorPicture, configs.AdminRole))).Methods("PUT")
	mux.Handle("/vendors/{id}/products", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.AddProductToVendor, configs.AdminRole))).Methods("POST")
	mux.Handle("/vendors/{vendorID}/products/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductOnVendor, configs.AdminRole))).Methods("PUT")
	mux.Handle("/vendors/{vendorID}/products/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.DeleteProductFromVendor, configs.AdminRole))).Methods("DELETE")
	mux.Handle("/vendors/{vendorID}/products/{id}/pictures", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductPicture, configs.AdminRole))).Methods("PUT")
	mux.Handle("/vendors/{id}/orders", csrfChecker.Check(orderDelivery.GetVendorOrders)).Methods("GET")
	mux.Handle("/vendors/{id}/reviews", csrfChecker.Check(reviewDelivery.GetVendorReviews)).Methods("GET")
	mux.Handle("/vendors/{vendorID}/orders/{id}", csrfChecker.Check(orderDelivery.UpdateOrderStatus)).Methods("PUT")
	mux.HandleFunc("/partners", partnerDelivery.Create).Methods("POST")
	mux.Handle("/partners/vendors", authChecker.Check(partnerDelivery.GetPartnerShops)).Methods("GET")
	mux.Handle("/carts", csrfChecker.Check(cartDelivery.AddToCart)).Methods("PUT")
	mux.Handle("/carts", csrfChecker.Check(cartDelivery.RemoveFromCart)).Methods("DELETE")
	mux.Handle("/carts", csrfChecker.Check(cartDelivery.GetCart)).Methods("GET")
	mux.Handle("/orders", csrfChecker.Check(orderDelivery.AddOrder)).Methods("POST")
	mux.Handle("/orders", csrfChecker.Check(orderDelivery.GetUserOrders)).Methods("GET")
	mux.Handle("/orders/{id}", csrfChecker.Check(orderDelivery.GetOrder)).Methods("GET")
	mux.Handle("/csrf", authChecker.Check(csrfDelivery.SetCSRF)).Methods("GET")
	mux.Handle("/reviews", csrfChecker.Check(reviewDelivery.AddReview)).Methods("POST")
	mux.Handle("/reviews", csrfChecker.Check(reviewDelivery.GetUserReviews)).Methods("GET")

	accessLogHandler := middleware.AccessLog(mux)
	corsHandler := middleware.CORS(accessLogHandler)
	siteHandler := middleware.Panic(corsHandler)

	logrus.Info("starting server at port ", configs.Port)
	logrus.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
