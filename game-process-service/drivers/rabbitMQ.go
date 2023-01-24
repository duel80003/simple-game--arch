package drivers

import (
	tools "github.com/duel80003/my-tools"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

var (
	RabbitMQConn   *amqp.Connection
	Exchange       string
	BetTableTMinus string
)

func RabbitMQInit() {
	var err error
	host := os.Getenv("RABBITMQ_ADDR")
	RabbitMQConn, err = amqp.Dial(host)
	if !InitExchange() {
		tools.Logger.Fatalf("rabbitMQ Init exchange failure")
	}
	if !InitRouters() {
		tools.Logger.Fatalf("rabbitMQ Init router failure")
	}
	if err != nil {
		tools.Logger.Fatalf("rabbitMQ connection failure: %s", err)
	}
	tools.Logger.Infof("rebbitMQ connected")
}

func InitExchange() bool {
	Exchange = os.Getenv("EXCHANGE")
	if Exchange == "" {
		return false
	}
	return true
}

func InitRouters() bool {
	BetTableTMinus = os.Getenv("BET_TABLE_T_MINUS")
	if BetTableTMinus == "" {
		tools.Logger.Errorf("empty key: %s", BetTableTMinus)
		return false
	}
	return true
}

func RabbitMQClose() {
	err := RabbitMQConn.Close()
	if err != nil {
		tools.Logger.Errorf("rabbitMQ disconnect failure, error: %s", err)
		return
	}
	tools.Logger.Infof("rabbitMQ discnnected")
}
