package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/friends/configs"
	cartDelivery "github.com/friends/internal/pkg/cart/delivery"
	cartRepo "github.com/friends/internal/pkg/cart/repository"
	cartUsecase "github.com/friends/internal/pkg/cart/usecase"
	chatDelivery "github.com/friends/internal/pkg/chat/delivery"
	chatRepository "github.com/friends/internal/pkg/chat/repository"
	chatUsecase "github.com/friends/internal/pkg/chat/usecase"
	"github.com/friends/internal/pkg/fileserver"
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
	"github.com/friends/internal/pkg/session"
	userDelivery "github.com/friends/internal/pkg/user/delivery"
	userRepo "github.com/friends/internal/pkg/user/repository"
	userUsecase "github.com/friends/internal/pkg/user/usecase"
	vendorDelivery "github.com/friends/internal/pkg/vendors/delivery"
	vendorRepo "github.com/friends/internal/pkg/vendors/repository"
	vendorUsecase "github.com/friends/internal/pkg/vendors/usecase"
	websocketpool "github.com/friends/internal/pkg/websocketPool"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StartAPIServer(dsn string) {
	logrus.SetLevel(logrus.DebugLevel)
	db, err := sql.Open(configs.Postgres, dsn)
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

	grpcSessionConn, err := grpc.Dial(
		"localhost"+configs.SessionServicePort,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}

	sessionClient := session.NewSessionWorkerClient(grpcSessionConn)

	grpcFileserverConn, err := grpc.Dial(
		"localhost"+configs.FileServerGRPCPort,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcSessionConn.Close()
	defer grpcFileserverConn.Close()

	fileserverClient := fileserver.NewUploadServiceClient(grpcFileserverConn)

	profRepo := profileRepo.NewProfileRepository(db)
	profUsecase := profileUsecase.NewProfileUsecase(profRepo, fileserverClient)
	profDelivery := profileDelivery.NewProfileDelivery(profUsecase)

	userHandler := userDelivery.NewUserHandler(userUsecase, sessionClient, profUsecase)

	vendRepo := vendorRepo.NewVendorRepository(db)
	vendUsecase := vendorUsecase.NewVendorUsecase(vendRepo, fileserverClient)
	vendDelivery := vendorDelivery.NewVendorDelivery(vendUsecase)

	cartRepo := cartRepo.NewCartRepository(db)
	cartUsecase := cartUsecase.NewCartUsecase(cartRepo, vendRepo)
	cartDelivery := cartDelivery.NewCartDelivery(cartUsecase)

	partnerDelivery := partnerDelivery.New(userUsecase, profUsecase, sessionClient, vendUsecase)

	wsPool := websocketpool.NewWebsocketPool()

	orderRepo := orderRepo.New(db)
	orderUsecase := orderUsecase.New(orderRepo, vendRepo)
	orderDelivery := orderDelivery.New(orderUsecase, vendUsecase, wsPool)

	reviewRepository := reviewRepository.New(db)
	reviewUsecase := reviewUsecase.New(reviewRepository, orderRepo, profRepo, vendRepo)
	reviewDelivery := reviewDelivery.New(reviewUsecase)

	chatRepository := chatRepository.New(db)
	chatUsecase := chatUsecase.New(chatRepository, profRepo, orderRepo)
	chatDelivery := chatDelivery.New(chatUsecase, orderUsecase, vendUsecase, wsPool)

	accessRighsChecker := middleware.NewAccessRightsChecker(userUsecase)

	authChecker := middleware.NewAuthChecker(sessionClient)
	csrfChecker := middleware.NewCSRFChecker(authChecker)

	mux := mux.NewRouter().PathPrefix(configs.APIURL).Subrouter()

	mux.HandleFunc("/users", userHandler.Create).Methods("POST")
	mux.Handle("/users", csrfChecker.Check(userHandler.Delete)).Methods("DELETE")

	mux.HandleFunc("/sessions", userHandler.Login).Methods("POST")
	mux.Handle("/sessions", csrfChecker.Check(userHandler.Logout)).Methods("DELETE")
	mux.HandleFunc("/sessions", userHandler.IsAuthorized).Methods("GET")

	mux.Handle("/profiles", csrfChecker.Check(profDelivery.Get)).Methods("GET")
	mux.Handle("/profiles", csrfChecker.Check(profDelivery.Update)).Methods("PUT")
	mux.Handle("/profiles/avatars", csrfChecker.Check(profDelivery.UpdateAvatar)).Methods("PUT")
	mux.Handle("/profiles/addresses", csrfChecker.Check(profDelivery.UpdateAddresses)).Methods("PUT")

	mux.HandleFunc("/vendors", vendDelivery.GetAll).Methods("GET")
	mux.HandleFunc("/vendors/nearest", vendDelivery.GetNearest).Methods("GET")
	mux.HandleFunc("/vendors/{id}", vendDelivery.GetVendor).Methods("GET")
	mux.Handle(
		"/vendors",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.CreateVendor, configs.AdminRole)),
	).Methods("POST")
	mux.Handle(
		"/vendors/{id}",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendor, configs.AdminRole)),
	).Methods("PUT")
	mux.Handle(
		"/vendors/{id}/pictures",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateVendorPicture, configs.AdminRole)),
	).Methods("PUT")
	mux.Handle(
		"/vendors/{id}/products",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.AddProductToVendor, configs.AdminRole)),
	).Methods("POST")
	mux.Handle(
		"/vendors/{vendorID}/products/{id}",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductOnVendor, configs.AdminRole)),
	).Methods("PUT")
	mux.Handle(
		"/vendors/{vendorID}/products/{id}",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.DeleteProductFromVendor, configs.AdminRole)),
	).Methods("DELETE")
	mux.Handle(
		"/vendors/{vendorID}/products/{id}/pictures",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(partnerDelivery.UpdateProductPicture, configs.AdminRole)),
	).Methods("PUT")
	mux.Handle(
		"/vendors/{id}/orders",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(orderDelivery.GetVendorOrders, configs.AdminRole)),
	).Methods("GET")
	mux.HandleFunc("/vendors/{id}/reviews", reviewDelivery.GetVendorReviews).Methods("GET")
	mux.Handle("/vendors/{vendorID}/orders/{id}", csrfChecker.Check(orderDelivery.UpdateOrderStatus)).Methods("PUT")
	mux.Handle(
		"/vendors/{id}/chats",
		csrfChecker.Check(accessRighsChecker.AccessRightsCheck(chatDelivery.GetVendorChats, configs.AdminRole)),
	).Methods("GET")
	mux.HandleFunc("/vendors/{id}/similar", vendDelivery.GetSimilar).Methods("GET")

	mux.HandleFunc("/partners", partnerDelivery.Create).Methods("POST")
	mux.Handle("/partners/vendors", csrfChecker.Check(partnerDelivery.GetPartnerShops)).Methods("GET")

	mux.Handle("/carts", csrfChecker.Check(cartDelivery.AddToCart)).Methods("PUT")
	mux.Handle("/carts", csrfChecker.Check(cartDelivery.RemoveFromCart)).Methods("DELETE")
	mux.Handle("/carts", csrfChecker.Check(cartDelivery.GetCart)).Methods("GET")

	mux.Handle("/orders", csrfChecker.Check(orderDelivery.AddOrder)).Methods("POST")
	mux.Handle("/orders", csrfChecker.Check(orderDelivery.GetUserOrders)).Methods("GET")
	mux.Handle("/orders/{id}", csrfChecker.Check(orderDelivery.GetOrder)).Methods("GET")

	mux.Handle("/reviews", csrfChecker.Check(reviewDelivery.AddReview)).Methods("POST")
	mux.Handle("/reviews", csrfChecker.Check(reviewDelivery.GetUserReviews)).Methods("GET")

	mux.Handle("/ws", authChecker.Check(chatDelivery.Upgrade)).Methods("GET")
	mux.Handle("/chats/{id}", csrfChecker.Check(chatDelivery.GetChat)).Methods("GET")

	mux.HandleFunc("/categories", vendDelivery.GetAllCategories).Methods("GET")

	accessLogHandler := middleware.AccessLog(mux)
	corsHandler := middleware.CORS(accessLogHandler)
	siteHandler := middleware.Panic(corsHandler)

	logrus.Info("starting server at port ", configs.Port)
	logrus.Fatal(http.ListenAndServe(configs.Port, siteHandler))
}
