package main

import (
	"context"
	tools "github.com/duel80003/my-tools"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"os/signal"
	"stress-testing/actions"
	"time"
)

func init() {
	err := godotenv.Load()
	tools.LogInit()
	if err != nil {
		tools.Logger.Fatalf("load env file failure: %s", err)
	}
	rand.Seed(time.Now().UnixMilli())
}

func main() {
	//addr := os.Getenv("WS_ADDR")
	//u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	//c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	actions.StartTesting()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)
	// Block until we receive our signal.
	<-c
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	actions.StopTesting()
}
