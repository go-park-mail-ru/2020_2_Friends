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

	profDelivery := profileDelivery.NewProfileDelivery(profUsecase)

	vendDelivery := vendorDelivery.NewVendorDelivery(vendUsecase)

	cartRepo := cartRepo.NewCartRepository(db)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepo, vendRepo)
	cartDelivery := cartDelivery.NewCartDelivery(cartUsecase)

	partnerDelivery := partnerDelivery.New(userUsecase, profUsecase, sessionUsecase, vendUsecase)

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

	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/users", userHandler.Create).Methods("POST")
	mux.HandleFunc("/api/v1/users", userHandler.Delete).Methods("DELETE")
	mux.HandleFunc("/api/v1/sessions", sessionDelivery.Create).Methods("POST")
	mux.HandleFunc("/api/v1/sessions", sessionDelivery.Delete).Methods("DELETE")
	mux.Handle("/api/v1/profiles", csrfChecker.Check(profDelivery.Get)).Methods("GET")
	mux.Handle("/api/v1/profiles", csrfChecker.Check(profDelivery.Update)).Methods("PUT")
	mux.Handle("/api/v1/profiles/avatars", csrfChecker.Check(profDelivery.UpdateAvatar)).Methods("PUT")
	mux.Handle("/api/v1/profiles/addresses", csrfChecker.Check(profDelivery.UpdateAddresses)).Methods("PUT")
	mux.HandleFunc("/api/v1/vendors", vendDelivery.GetAll).Methods("GET")
	mux.HandleFunc("/api/v1/vendors/{id}", vendDelivery.GetVendor).Methods("GET")
	mux.Handle("/api/v1/vendors", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.CreateVendor, configs.AdminRole))).Methods("POST")
	mux.Handle("/api/v1/vendors/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendor, configs.AdminRole))).Methods("PUT")
	mux.Handle("/api/v1/vendors/{id}/pictures", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendorPicture, configs.AdminRole))).Methods("PUT")
	mux.Handle("/api/v1/vendors/{id}/products", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.AddProductToVendor, configs.AdminRole))).Methods("POST")
	mux.Handle("/api/v1/vendors/{vendorID}/products/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductOnVendor, configs.AdminRole))).Methods("PUT")
	mux.Handle("/api/v1/vendors/{vendorID}/products/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.DeleteProductFromVendor, configs.AdminRole))).Methods("DELETE")
	mux.Handle("/api/v1/vendors/{vendorID}/products/{id}/pictures", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductPicture, configs.AdminRole))).Methods("PUT")
	mux.HandleFunc("/api/v1/partners", partnerDelivery.Create).Methods("POST")
	mux.Handle("/api/v1/partners/vendors", authChecker.Check(partnerDelivery.GetPartnerShops)).Methods("GET")
	mux.Handle("/api/v1/carts", csrfChecker.Check(cartDelivery.AddToCart)).Methods("PUT")
	mux.Handle("/api/v1/carts", csrfChecker.Check(cartDelivery.RemoveFromCart)).Methods("DELETE")
	mux.Handle("/api/v1/carts", csrfChecker.Check(cartDelivery.GetCart)).Methods("GET")
	mux.Handle("/api/v1/csrf", authChecker.Check(csrfDelivery.SetCSRF)).Methods("GET")

	mux.HandleFunc("/api/v2/users", userHandler.Create2).Methods("POST")
	mux.HandleFunc("/api/v2/users", userHandler.Delete).Methods("DELETE")
	mux.HandleFunc("/api/v2/sessions", sessionDelivery.Create2).Methods("POST")
	mux.HandleFunc("/api/v2/sessions", sessionDelivery.Delete).Methods("DELETE")
	mux.Handle("/api/v2/profiles", csrfChecker.Check(profDelivery.Get)).Methods("GET")
	mux.Handle("/api/v2/profiles", csrfChecker.Check(profDelivery.Update)).Methods("PUT")
	mux.Handle("/api/v2/profiles/avatars", csrfChecker.Check(profDelivery.UpdateAvatar)).Methods("PUT")
	mux.Handle("/api/v2/profiles/addresses", csrfChecker.Check(profDelivery.UpdateAddresses)).Methods("PUT")
	mux.HandleFunc("/api/v2/vendors", vendDelivery.GetAll).Methods("GET")
	mux.HandleFunc("/api/v2/vendors/{id}", vendDelivery.GetVendor).Methods("GET")
	mux.Handle("/api/v2/vendors", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.CreateVendor, configs.AdminRole))).Methods("POST")
	mux.Handle("/api/v2/vendors/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendor, configs.AdminRole))).Methods("PUT")
	mux.Handle("/api/v2/vendors/{id}/pictures", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendorPicture, configs.AdminRole))).Methods("PUT")
	mux.Handle("/api/v2/vendors/{id}/products", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.AddProductToVendor, configs.AdminRole))).Methods("POST")
	mux.Handle("/api/v2/vendors/{vendorID}/products/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductOnVendor, configs.AdminRole))).Methods("PUT")
	mux.Handle("/api/v2/vendors/{vendorID}/products/{id}", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.DeleteProductFromVendor, configs.AdminRole))).Methods("DELETE")
	mux.Handle("/api/v2/vendors/{vendorID}/products/{id}/pictures", csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductPicture, configs.AdminRole))).Methods("PUT")
	mux.HandleFunc("/api/v2/partners", partnerDelivery.Create2).Methods("POST")
	mux.Handle("/api/v2/partners/vendors", authChecker.Check(partnerDelivery.GetPartnerShops)).Methods("GET")
	mux.Handle("/api/v2/carts", csrfChecker.Check(cartDelivery.AddToCart)).Methods("PUT")
	mux.Handle("/api/v2/carts", csrfChecker.Check(cartDelivery.RemoveFromCart)).Methods("DELETE")
	mux.Handle("/api/v2/carts", csrfChecker.Check(cartDelivery.GetCart)).Methods("GET")
	mux.Handle("/api/v2/csrf", authChecker.Check(csrfDelivery.SetCSRF)).Methods("GET")

	accessLogHandler := middleware.AccessLog(mux)
	corsHandler := middleware.CORS(accessLogHandler)
	siteHandler := middleware.Panic(corsHandler)

	logrus.Info("starting server at port ", configs.Port)
	logrus.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
