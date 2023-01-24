package state

import (
	tools "github.com/duel80003/my-tools"
	"github.com/joho/godotenv"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("../.env")
	tools.LogInit()
	if err != nil {
		tools.Logger.Panic("load env file failure")
	}
}

func TestState(t *testing.T) {
	StartStateMachine()
	time.Sleep(40 * time.Second)
}
