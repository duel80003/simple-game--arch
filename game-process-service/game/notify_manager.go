package game

import (
	"game-process-service/models"
	"sync"
	"time"
)

var (
	notifyManager     *NotifyManager
	notifyManagerOnce sync.Once
)

type NotifyManager struct {
	Channels []chan *models.NotificationEvent
}

func initNotifyManager() {
	notifyManagerOnce.Do(func() {
		notifyManager = new(NotifyManager)
		for i := 0; i < 10; i++ {
			ch := make(chan *models.NotificationEvent, 1)
			notifyManager.Channels = append(notifyManager.Channels, ch)
			chanHandler := newChanHandler(ch)
			go chanHandler.startProcess()
		}
	})
}

func GetNotifyManager() *NotifyManager {

	return notifyManager
}

func (manager *NotifyManager) Notify(event *models.NotificationEvent) {
	for _, v := range manager.Channels {
		v <- event
	}
}

func (manager *NotifyManager) Join(event *models.NotificationEvent) {
	now := time.Now().UnixMilli()
	index := now % 10
	manager.Channels[index] <- event
}

func (manager *NotifyManager) Leave(event *models.NotificationEvent) {
	for _, v := range manager.Channels {
		v <- event
	}
}
