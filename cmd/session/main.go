package main

import (
	"context"
	"log"
	"net"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/session"
	delivery "github.com/friends/internal/pkg/session/delivery"
	repository "github.com/friends/internal/pkg/session/repository"
	usecase "github.com/friends/internal/pkg/session/usecase"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StartSessionService() {
	lis, err := net.Listen("tcp", configs.SessionServicePort)
	if err != nil {
		log.Fatalln("can't start session service: ", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     configs.RedisAddr,
		Password: "",
		DB:       0,
	})

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("redis not available: ", err)
	}

	sessionRepo := repository.NewSessionRedisRepo(redisClient)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo)

	server := grpc.NewServer()

	sessionDelivery := delivery.NewSessionDelivery(sessionUsecase)

	session.RegisterSessionWorkerServer(server, sessionDelivery)

	logrus.Info("starting session service at port ", configs.SessionServicePort)
	log.Fatal(server.Serve(lis))
}

func main() {
	StartSessionService()
}
