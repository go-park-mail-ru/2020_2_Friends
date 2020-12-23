package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/fileserver"
	"github.com/friends/internal/pkg/fileserver/delivery"
	"github.com/friends/internal/pkg/fileserver/repository"
	"github.com/friends/internal/pkg/fileserver/usecase"
	"github.com/friends/internal/pkg/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StartFileServer() {
	mux := http.NewServeMux()

	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir(os.Getenv("img_path")),
	)
	mux.Handle("/data/", staticHandler)

	corsHandler := middleware.CORS(mux)
	siteHandler := middleware.Panic(corsHandler)

	logrus.Info("starting fileserver at port ", configs.FileServerPort)
	logrus.Fatal(http.ListenAndServe(configs.FileServerPort, siteHandler))
}

func StartGRPCServer() {
	repository := repository.New()
	usecase := usecase.New(repository)

	delivery := delivery.New(usecase)

	server := grpc.NewServer()

	fileserver.RegisterUploadServiceServer(server, delivery)

	lis, err := net.Listen("tcp", configs.FileServerGRPCPort)
	if err != nil {
		log.Fatalln("can't start session service: ", err)
	}

	logrus.Info("starting fileserver service at port ", configs.FileServerGRPCPort)
	log.Fatal(server.Serve(lis))
}

func main() {
	go StartFileServer()
	StartGRPCServer()
}
