package main

import (
	"context"
	"fmt"
	"game-process-service/drivers"
	"game-process-service/game"
	"game-process-service/servers"
	tools "github.com/duel80003/my-tools"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func init() {
	err := godotenv.Load()
	tools.LogInit()
	if err != nil {
		tools.Logger.Infof("load env file failure")
	}

	drivers.InitStateGrpcConn()
	drivers.RedisInit()
	drivers.RabbitMQInit()
}

func serverStart() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		tools.Logger.Panicf("empty grpc port")
	}
	s := servers.GRpcServers()
	tools.Logger.Infof("server start: %s", grpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	game.InitGameRoom()
	time.Sleep(1 * time.Second)
	go game.StartGameFlow()
	go serverStart()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)
	// Block until we receive our signal.
	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	drivers.CloseStateGrpcConn()
	drivers.RedisFlushAll()
	drivers.RedisClose()
	drivers.RabbitMQClose()
}
