package workers

import (
	"encoding/json"
	tools "github.com/duel80003/my-tools"
	"github.com/rabbitmq/amqp091-go"
	. "main-service/drivers"
	"main-service/handler"
	"main-service/models"
)

func StateWorkerStart() {
	go func() {
		ch, msgs := workerInit(Exchange, TableState)
		defer ch.Close()
		for d := range msgs {
			stateMsgHandler(d)
		}
	}()
}

func stateMsgHandler(d amqp091.Delivery) {
	tools.Logger.Infof("[%s] receives a message: %s", TableState, d.Body)
	state := new(models.StateInfo)
	err := json.Unmarshal(d.Body, state)
	if err != nil {
		tools.Logger.Errorf("unmarshal error: %s", err)
		return
	}
	handler.GetWsHandler().Broadcast(handler.State, state)
}
