package game

import (
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
	mqRepo  repositories.MqRepository
}

func newChanHandler(ch chan *models.NotificationEvent) *ChanHandler {
	chanHandler := new(ChanHandler)
	chanHandler.ch = ch
	chanHandler.players = make(map[string]string)
	chanHandler.ticker = time.NewTicker(1 * time.Second)
	chanHandler.mqRepo = repositories.NewMqRepository()
	return chanHandler
}

func (handler *ChanHandler) startProcess() {
	for info := range handler.ch {
		switch info.Type {
		case models.NotifyState:
			handler.stateNotify()
		case models.NotifyBetZoneInfo:
			handler.BetZoneInfos(info.TMinus)
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
	handler.mux.RLock()
	defer handler.mux.RUnlock()
	for key, value := range handler.players {
		event := new(models.Event)
		event.Exchange = ExchangeGameState
		event.Router = TableState
		state := new(models.StateInfo)
		state.State = GetRoom().State
		event.Data = &models.EventData{
			Data: state,
		}
		event.Data.Session = key
		event.Data.PlayerID = value
		//repositories.PublishEvent(context.TODO(), event)
		handler.mqRepo.PublishEvent(event)
	}
}

func (handler *ChanHandler) BetZoneInfos(tMinus int32) {
	defer handler.ticker.Stop()
	handler.ticker.Reset(1 * time.Second)
	//tools.Logger.Infof("map count: %d", len(handler.players))
	for {
		select {
		case <-handler.ticker.C:
			betZones := repositories.GetRoomBetInfo(GetRoom().RoomID)
			betZoneInfos := new(models.BetZoneInfos)
			betZoneInfos.BetZones = betZones
			betZoneInfos.TMinus = tMinus
			handler.mux.RLock()
			for key, value := range handler.players {
				event := new(models.Event)
				event.Exchange = ExchangeBetInfo
				event.Router = BetTableTMinus
				event.Data = &models.EventData{
					Data: betZoneInfos,
				}
				event.Data.Session = key
				event.Data.PlayerID = value
				//repositories.PublishEvent(context.TODO(), event)
				handler.mqRepo.PublishEvent(event)
			}
			handler.mux.RUnlock()
			tMinus--
			if tMinus <= 0 {
				return
			}
		}
	}
}
