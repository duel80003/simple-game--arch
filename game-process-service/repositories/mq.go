package repositories

import (
	"context"
	"encoding/json"
	. "game-process-service/drivers"
	"game-process-service/models"
	tools "github.com/duel80003/my-tools"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishEvent(ctx context.Context, event *models.Event) (err error) {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		tools.Logger.Errorf("rabbitMQ get channel error: %s", err)
		return
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		event.Exchange, // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		tools.Logger.Errorf("rabbitMQ ExchangeDeclare error: %s", err)
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
