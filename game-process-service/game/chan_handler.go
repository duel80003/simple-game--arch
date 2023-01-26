package game

import (
	"context"
	. "game-process-service/drivers"
	"game-process-service/models"
	"game-process-service/repositories"
	"sync"
	"time"
)

type ChanHandler struct {
	players map[string]string
	ch      chan *models.NotificationEvent
	mux     sync.RWMutex
	ticker  *time.Ticker
}

func (handler *ChanHandler) startProcess() {
	for info := range handler.ch {
		switch info.Type {
		case models.NotifyState:
			handler.stateNotify()
		case models.NotifyBetZoneInfo:
			handler.BetZoneInfos(info.TMinus, info.BetZones)
		case models.PlayerJoin:
			handler.playerJoin(info.SID, info.PID)
		case models.PlayerLeave:
			handler.playerLeave(info.SID)
		}
	}
}

func (handler *ChanHandler) playerJoin(sid, pid string) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	handler.players[sid] = pid
}

func (handler *ChanHandler) playerLeave(sid string) {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	delete(handler.players, sid)
}

func (handler *ChanHandler) stateNotify() {
	event := new(models.Event)
	event.Exchange = ExchangeGameState
	event.Router = TableState
	state := new(models.StateInfo)
	state.State = GetRoom().State
	event.Data = &models.EventData{
		Data: state,
	}
	handler.mux.RLock()
	defer handler.mux.RUnlock()
	for key, value := range handler.players {
		event.Data.Session = key
		event.Data.PlayerID = value
		repositories.PublishEvent(context.TODO(), event)
	}
}

func (handler *ChanHandler) BetZoneInfos(tMinus int32, betZones *models.BetZones) {
	defer handler.ticker.Stop()
	handler.ticker.Reset(1 * time.Second)
	//tools.Logger.Infof("map count: %d", len(handler.players))
	for {
		select {
		case <-handler.ticker.C:
			betZoneInfos := new(models.BetZoneInfos)
			betZoneInfos.BetZones = betZones
			betZoneInfos.TMinus = tMinus
			event := new(models.Event)
			event.Exchange = ExchangeBetInfo
			event.Router = BetTableTMinus
			event.Data = &models.EventData{
				Data: betZoneInfos,
			}
			handler.mux.RLock()
			for key, value := range handler.players {
				event.Data.Session = key
				event.Data.PlayerID = value
				repositories.PublishEvent(context.TODO(), event)
			}
			handler.mux.RUnlock()
			tMinus--
			if tMinus <= 0 {
				return
			}
		}
	}
}
