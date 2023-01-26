package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	. "game-process-service/drivers"
	"game-process-service/models"
	tools "github.com/duel80003/my-tools"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishEvent(ctx context.Context, event *models.Event) (err error) {
	ch := GetChannel(event.Exchange)
	if ch == nil {
		tools.Logger.Errorf("[PublishEvent] channel not exist: %s", event.Exchange)
		err = fmt.Errorf("emtpy channel")
		return
	}
	body, err := json.Marshal(event.Data)
	err = ch.PublishWithContext(
		ctx,
		event.Exchange, // exchange
		event.Router,   // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	return err
}
