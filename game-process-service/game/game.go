package game

import (
	"context"
	"fmt"
	"game-process-service/config"
	. "game-process-service/drivers"
	"game-process-service/models"
	proto "game-process-service/proto/gen/v1"
	"game-process-service/repositories"
	tools "github.com/duel80003/my-tools"
	"github.com/google/uuid"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var room *Room

func GetRoom() *Room {
	return room
}

type Room struct {
	RoomID        string
	RoomType      string
	GameID        string
	Chips         []int32
	PlayerSession map[string]string
	State         proto.State
	mux           sync.RWMutex
	ticker        *time.Ticker
}

func InitGameRoom() {
	roomType := os.Getenv("ROOM_TYPE")
	chips := os.Getenv("CHIPS")
	gid := os.Getenv("GAME_ID")
	room = new(Room)
	room.GameID = gid
	room.RoomID = fmt.Sprintf("room_%s", uuid.NewString())
	room.RoomType = roomType
	room.Chips = parseChips(chips)
	room.PlayerSession = make(map[string]string)
	room.ticker = time.NewTicker(1 * time.Second)
}

func parseChips(chipStr string) (chips []int32) {
	str := strings.Split(chipStr, ":")
	for _, s := range str {
		i, err := strconv.Atoi(s)
		if err != nil {
			tools.Logger.Fatalf("invalid chip info: %s", s)
		}
		chips = append(chips, int32(i))
	}
	return
}

func (room *Room) Reset() {
	repositories.ResetRoomBetInfo(room.RoomID, []string{config.BetZoneMap[0], config.BetZoneMap[1]})
	sids := make([]string, 0, len(room.PlayerSession))
	for sid := range room.PlayerSession {
		sids = append(sids, sid)
	}
	repositories.ResetPlayerBetInfo(sids, []string{config.BetZoneMap[0], config.BetZoneMap[1]})
}

func (room *Room) Join(pid, sid string) {
	room.mux.Lock()
	defer room.mux.Unlock()
	room.PlayerSession[sid] = pid
	repositories.SetPlayer(sid, pid, room.GameID)
}

func (room *Room) Leave(sid string) {
	room.mux.Lock()
	defer room.mux.Unlock()
	delete(room.PlayerSession, sid)
	repositories.RemovePlayer(sid)
}

func (room *Room) Bet(sid string, betZone, bet int32) bool {
	tools.Logger.Infof("[Bet] current state: %v", room.State)
	if room.State != proto.State_STATE_START_BET {
		return false
	}
	room.mux.RLock()
	if _, ok := room.PlayerSession[sid]; !ok {
		return false
	}
	room.mux.RUnlock()
	repositories.UpdateRoomBetInfo(room.RoomID, config.BetZoneMap[betZone], bet)
	repositories.UpdatePlayerBetInfo(sid, config.BetZoneMap[betZone], bet)
	return true
}

func (room *Room) BetZoneInfos(tMinus int32) {
	defer room.ticker.Stop()
	room.ticker.Reset(1 * time.Second)
	for {
		select {
		case <-room.ticker.C:
			betZones := repositories.GetRoomBetInfo(room.RoomID)
			betZoneInfos := new(models.BetZoneInfos)
			betZoneInfos.BetZones = betZones
			betZoneInfos.TMinus = tMinus
			betZoneInfos.GameID = room.GameID
			event := new(models.Event)
			event.Exchange = ExchangeBetInfo
			event.Router = BetTableTMinus
			event.Data = betZoneInfos
			repositories.PublishEvent(context.TODO(), event)
			tMinus--
			if tMinus <= 0 {
				return
			}
		}
	}
}

func (room *Room) StateNotify() {
	event := new(models.Event)
	event.Exchange = ExchangeGameState
	event.Router = TableState
	state := new(models.StateInfo)
	state.State = room.State
	state.GameID = room.GameID
	event.Data = state
	repositories.PublishEvent(context.TODO(), event)
}
