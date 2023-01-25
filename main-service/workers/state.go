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
		ch, msgs := workerInit(ExchangeGameState, TableState)
		defer ch.Close()
		for d := range msgs {
			stateMsgHandler(d)
		}
	}()
}

func stateMsgHandler(d amqp091.Delivery) {
	tools.Logger.Debugf("[%s] receives a message: %s", TableState, d.Body)
	eventData := new(models.EventData)
	err := json.Unmarshal(d.Body, eventData)
	if err != nil {
		tools.Logger.Errorf("unmarshal error: %s", err)
		return
	}
	handler.GetWsHandler().Broadcast(handler.State, eventData)
}
