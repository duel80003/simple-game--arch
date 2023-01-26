package repositories

import (
	"context"
	"fmt"
	. "game-process-service/drivers"
	"game-process-service/models"
	"sync"
	"testing"
	"time"
)

func TestPushEvent(t *testing.T) {
	event := &models.Event{
		Exchange: ExchangeBetInfo,
		Router:   BetTableTMinus,
		Data: &models.EventData{
			Session:  "test",
			PlayerID: "test_player",
			Data:     "data",
		},
	}
	err := PublishEvent(context.TODO(), event)
	if err != nil {
		t.Errorf("publish event error")
		return
	}
	time.Sleep(1 * time.Second)
}

func TestPushEvents(t *testing.T) {
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		func(index int) {
			event := &models.Event{
				Exchange: ExchangeBetInfo,
				Router:   BetTableTMinus,
				Data: &models.EventData{
					Session:  fmt.Sprintf("sid_%d", index+1),
					PlayerID: fmt.Sprintf("pid_%d", index+1),
					Data:     "data",
				},
			}
			t.Logf("push event: %+v", event)
			err := PublishEvent(context.TODO(), event)
			wg.Done()
			if err != nil {
				t.Errorf("publish event error")
				return
			}
		}(i)
	}
	wg.Wait()
	time.Sleep(5 * time.Second)
}

func TestPushEvents2(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	times := 0
	for {
		select {
		case <-ticker.C:
			count := 1000
			var wg sync.WaitGroup
			wg.Add(count)
			for i := 0; i < count; i++ {
				func(index int) {
					event := &models.Event{
						Exchange: ExchangeBetInfo,
						Router:   BetTableTMinus,
						Data: &models.EventData{
							Session:  fmt.Sprintf("sid_%d", index+1),
							PlayerID: fmt.Sprintf("pid_%d", index+1),
							Data:     "data",
						},
					}
					t.Logf("push event: %+v", event)
					err := PublishEvent(context.TODO(), event)
					wg.Done()
					if err != nil {
						t.Errorf("publish event error")
						return
					}
				}(i)
			}
			wg.Wait()
			times++
			if times == 20 {
				return
			}
		}
	}
}
