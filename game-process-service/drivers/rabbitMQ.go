package drivers

import (
	tools "github.com/duel80003/my-tools"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"sync"
)

var (
	RabbitMQConn      *amqp.Connection
	ChannelBetInfo    *amqp.Channel
	ChannelGameState  *amqp.Channel
	ExchangeBetInfo   string
	ExchangeGameState string
	BetTableTMinus    string
	TableState        string
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
	ExchangeBetInfo = os.Getenv("EXCHANGE_BET_INFO")
	if ExchangeBetInfo == "" {
		return false
	}
	ExchangeGameState = os.Getenv("EXCHANGE_STATE")
	if ExchangeGameState == "" {
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
	TableState = os.Getenv("TABLE_STATE")
	if TableState == "" {
		return false
	}
	return true
}

func RabbitMQClose() {
	ChannelBetInfo.Close()

	err := RabbitMQConn.Close()
	if err != nil {
		tools.Logger.Errorf("rabbitMQ disconnect failure, error: %s", err)
		return
	}
	tools.Logger.Infof("rabbitMQ discnnected")
}

func GetChannel(exchange string) *amqp.Channel {
	switch exchange {
	case ExchangeBetInfo:
		return getChannelBetInfo()
	case ExchangeGameState:
		return getChannelState()
	default:
		return nil
	}
}

func getChannelBetInfo() *amqp.Channel {
	if RabbitMQConn.IsClosed() {
		RabbitMQConn.Close()
		RabbitMQInit()
	}

	if ChannelBetInfo == nil || ChannelBetInfo.IsClosed() {
		var mux sync.Mutex
		mux.Lock()
		ch, err := RabbitMQConn.Channel()
		if err != nil {
			tools.Logger.Errorf("rabbitMQ get getChannelBetInfo error: %s", err)
		}
		err = ch.ExchangeDeclare(
			ExchangeBetInfo, // name
			"fanout",        // type
			true,            // durable
			false,           // auto-deleted
			false,           // internal
			false,           // no-wait
			nil,             // arguments
		)
		if err != nil {
			tools.Logger.Errorf("rabbitMQ ExchangeDeclare error: %s", err)
		}
		ChannelBetInfo = ch
		mux.Unlock()
	}
	return ChannelBetInfo
}

func getChannelState() *amqp.Channel {
	if RabbitMQConn.IsClosed() {
		RabbitMQConn.Close()
		RabbitMQInit()
	}

	if ChannelGameState == nil || ChannelGameState.IsClosed() {
		var mux sync.Mutex
		mux.Lock()
		ch, err := RabbitMQConn.Channel()
		if err != nil {
			tools.Logger.Errorf("rabbitMQ get ChannelGameState error: %s", err)
		}
		err = ch.ExchangeDeclare(
			ExchangeGameState, // name
			"fanout",          // type
			true,              // durable
			false,             // auto-deleted
			false,             // internal
			false,             // no-wait
			nil,               // arguments
		)
		if err != nil {
			tools.Logger.Errorf("rabbitMQ ExchangeDeclare error: %s", err)
		}
		ChannelGameState = ch
		mux.Unlock()
	}
	return ChannelGameState
}

func InitChannels() {
	getChannelBetInfo()
	getChannelState()
}
