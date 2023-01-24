package workers

import (
	"encoding/json"
	tools "github.com/duel80003/my-tools"
	"github.com/rabbitmq/amqp091-go"
	. "main-service/drivers"
	"main-service/handler"
	"main-service/models"
)

func betZoneInfoWorkerStart() {
	go func() {
		ch, msgs := workerInit(ExchangeBetInfo, BetTableTMinus)
		defer ch.Close()
		for d := range msgs {
			betZoneInfoMsgHandler(d)
		}
	}()
}

func betZoneInfoMsgHandler(d amqp091.Delivery) {
	tools.Logger.Debugf("[%s] receives a message: %s", BetTableTMinus, d.Body)
	betZoneInfo := &models.BetZoneInfos{}
	err := json.Unmarshal(d.Body, betZoneInfo)
	if err != nil {
		tools.Logger.Errorf("unmarshal error: %s", err)
		return
	}
	handler.GetWsHandler().Broadcast(handler.BetZoneInfos, betZoneInfo.GameID, betZoneInfo)
}
