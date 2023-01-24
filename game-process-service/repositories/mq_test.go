package repositories

import (
	"context"
	. "game-process-service/drivers"
	"game-process-service/models"
	"testing"
	"time"
)

func TestPushEvent(t *testing.T) {
	event := &models.Event{
		Exchange: ExchangeBetInfo,
		Router:   BetTableTMinus,
		Data:     "test",
	}
	err := PublishEvent(context.TODO(), event)
	if err != nil {
		t.Errorf("publish event error")
		return
	}
	time.Sleep(1 * time.Second)
}
