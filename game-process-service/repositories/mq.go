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

type MqRepository interface {
	PublishEvent(event *models.Event)
}

type RbMQ struct {
	ch      *amqp.Channel
	eventCh chan *models.Event
}

func NewMqRepository() MqRepository {
	conn, _ := GetRabbitMQConn()
	ch, err := conn.Channel()
	if err != nil {
		tools.Logger.Errorf("[NewMqRepository] init ch error: %s", err)
		return nil
	}
	eventCh := make(chan *models.Event, 200)

	instance := &RbMQ{
		ch:      ch,
		eventCh: eventCh,
	}
	go instance.handler()
	return instance
}

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
	if err != nil {
		tools.Logger.Errorf("[PublishEvent] PublishWithContext error: %s", err)
	}
	return err
}

func (mq *RbMQ) PublishEvent(event *models.Event) {
	mq.eventCh <- event
}

func (mq *RbMQ) publishEvent(ctx context.Context, event *models.Event) (err error) {
	err = mq.ch.ExchangeDeclare(
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
	}
	body, err := json.Marshal(event.Data)
	err = mq.ch.PublishWithContext(
		ctx,
		event.Exchange, // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		tools.Logger.Errorf("[PublishEvent] PublishWithContext error: %s, ch: %v", err, &mq.ch)
	}
	return err
}

func (mq *RbMQ) handler() {
	for event := range mq.eventCh {
		mq.publishEvent(context.TODO(), event)
	}
}
