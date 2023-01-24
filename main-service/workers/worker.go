package workers

import (
	tools "github.com/duel80003/my-tools"
	"github.com/rabbitmq/amqp091-go"
	. "main-service/drivers"
)

func failOnError(err error, msg string) {
	if err != nil {
		tools.Logger.Panicf("rabbit mq %s, error: %s", msg, err)
	}
}

func workerInit(exchange, routerKey string) (*amqp091.Channel, <-chan amqp091.Delivery) {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		tools.Logger.Errorf("rabbitMQ get channel error: %s", err)
	}
	failOnError(err, "open channel")

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "declare exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "declare queue")

	err = ch.QueueBind(
		q.Name,    // queue name
		routerKey, // routing key
		exchange,  // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "consume queue")
	return ch, msgs

}

func StartWorkers() {
	betZoneInfoWorkerStart()
}
