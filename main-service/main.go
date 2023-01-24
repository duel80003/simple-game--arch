package main

import (
	"context"
	tools "github.com/duel80003/my-tools"
	"github.com/gobwas/ws"
	"github.com/joho/godotenv"
	"main-service/drivers"
	"main-service/handler"
	"main-service/workers"
	"net/http"
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
	drivers.InitGameGrpcConn()
	drivers.RedisInit()
	drivers.RabbitMQInit()
}

func main() {
	addr := os.Getenv("WS_PORT")
	go func() {
		wsHandler := handler.GetWsHandler()
		http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				tools.Logger.Errorf("conn error: %s", err)
				return
			}
			wsHandler.Run(conn)
		}))
	}()
	workers.StartWorkers()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)
	// Block until we receive our signal.
	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	drivers.CloseGameGrpcConn()
	drivers.RedisClose()
	drivers.RabbitMQClose()
}
