package workers

import (
	tools "github.com/duel80003/my-tools"
	"github.com/rabbitmq/amqp091-go"
	. "main-service/drivers"
)

func betZoneInfoWorkerStart() {
	go func() {
		ch, msgs := workerInit(Exchange, BetTableTMinus)
		defer ch.Close()
		for d := range msgs {
			userMsgHandler(d)
		}
	}()
}

func userMsgHandler(d amqp091.Delivery) {
	tools.Logger.Infof("[%s] receives a message: %s", BetTableTMinus, d.Body)
	//user := &models.User{}
	//err := json.Unmarshal(d.Body, user)
	//if err != nil {
	//	logger.Errorf("unmarshal error: %s", err)
	//	return
	//}
	//err = validate.Struct(user)
	//if err != nil {
	//	logger.Errorf("invalid data %s", err)
	//	return
	//}
	//now := time.Now().Unix()
	//user.CreatedAt = now
	//user.UpdatedAt = now
	//user.Rooms = make([]string, 0)
	//err = GetUserRepository().Insert(context.TODO(), user)
	//if err != nil {
	//	// TODO put message back to mq
	//	logger.Errorf("handle message error: %s", err)
	//}
}
