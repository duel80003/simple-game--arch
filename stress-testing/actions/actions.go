package actions

import (
	"encoding/json"
	"fmt"
	tools "github.com/duel80003/my-tools"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"stress-testing/models"
	"time"
)

const (
	PID = "demo_player_00001"
)

var (
	players  []*models.Player
	chips    []int32
	betZones []int32
	betTimes int
)

func initPlayerList() {
	countStr := os.Getenv("PLAYER_COUNT")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		tools.Logger.Fatalf("invalid count")
	}
	players = make([]*models.Player, 0, count)
	for i := 0; i < count; i++ {
		player := new(models.Player)
		player.PlayerID = fmt.Sprintf("demo_player_%05d", i+1)
		players = append(players, player)
	}
}

func initConn() {
	addr := os.Getenv("WS_ADDR")
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	for i := range players {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			tools.Logger.Fatalf("init ws connection error: %s", err)
		}
		players[i].Conn = c
		time.Sleep(time.Millisecond * 30)
	}
}

func StartTesting() {
	initPlayerList()
	tools.Logger.Infof("players count: %d", len(players))
	initConn()
	var err error
	betTimesStr := os.Getenv("BET_TIMES")
	betTimes, err = strconv.Atoi(betTimesStr)
	if err != nil {
		tools.Logger.Fatalf("invalid bet times")
	}
	tools.Logger.Infof("bet times: %d", betTimes)
	for i := range players {
		go func(p *models.Player) {
			messageHandler(p)
		}(players[i])
	}
}

func StopTesting() {
	tools.Logger.Infof("stopping....")
	for i := range players {
		go leaveAction(players[i])
	}

	var basicCount float64 = 5
	countStr := os.Getenv("PLAYER_COUNT")
	count, _ := strconv.ParseFloat(countStr, 64)
	count /= 20
    tools.Logger.Infof("wiat: %.2f seconds", count)
	time.Sleep(time.Second * time.Duration(count+basicCount))
}

func messageHandler(player *models.Player) {
	go func() {
		for {
			_, message, err := player.Conn.ReadMessage()
			if err != nil {
				tools.Logger.Errorf("read: %s", err)
				return
			}
			if player.PlayerID == PID {
				tools.Logger.Infof("recv: %s", message)
			}
			var res models.Response
			err = json.Unmarshal(message, &res)
			if err != nil {
				tools.Logger.Fatalf("res unmarshal error: %s", err)
			}
			switch res.Topic {
			case Join:
				tools.Logger.Infof("player: %s join game success", player.PlayerID)
				handleJoinMsg(&res)
			case Leave:
				tools.Logger.Infof("player: %s leave game success", player.PlayerID)
				return
			case State:
				stateI := res.Data[State]
				state, ok := stateI.(float64)
				if !ok {
					continue
				}
				if state == 2 {
					betAction(player)
				}
			}
		}
	}()
	go joinAction(player)
}

func joinAction(player *models.Player) {
	joinReq := new(models.Request)
	joinReq.Topic = Join
	joinReq.Data = map[string]interface{}{"playerId": player.PlayerID}
	err := player.Conn.WriteJSON(joinReq)
	if err != nil {
		tools.Logger.Errorf("player: %s join action error: %s", player.PlayerID, err)
	}
}

func handleJoinMsg(res *models.Response) {
	if len(chips) != 0 {
		return
	}
	parseChips(res)
	parseBetZones(res)
}

func leaveAction(player *models.Player) {
	leaveReq := new(models.Request)
	leaveReq.Topic = Leave
	leaveReq.Data = map[string]interface{}{"playerId": player.PlayerID}
	err := player.Conn.WriteJSON(leaveReq)
	if err != nil {
		tools.Logger.Errorf("player: %s leave action error: %s", player.PlayerID, err)
	}
}

func betAction(player *models.Player) {
	for i := 0; i < betTimes; i++ {
		time.AfterFunc(getRandBetTime(), func() {
			if player.PlayerID == PID {
				tools.Logger.Infof("player: %s do bet", player.PlayerID)
			}
			player.BetMux.Lock()
			defer player.BetMux.Unlock()
			betReq := new(models.Request)
			betReq.Topic = Bet
			betReq.Data = map[string]interface{}{
				"chip":    getRandChip(),
				"betZone": getRandBetZone(),
			}
			err := player.Conn.WriteJSON(betReq)
			if err != nil {
				tools.Logger.Errorf("player: %s leave action error: %s", player.PlayerID, err)
			}
		})
	}
}

func getRandChip() int32 {
	max := len(chips)
	return chips[rand.Intn(max)]
}

func getRandBetZone() int32 {
	max := len(betZones)
	return betZones[rand.Intn(max)]
}

func getRandBetTime() time.Duration {
	min := 1
	max := 8
	return time.Second * time.Duration(rand.Intn(max-min)+min)
}

func parseChips(res *models.Response) {
	chipsI, ok := res.Data["chips"]
	if !ok {
		tools.Logger.Fatal("chips is required")
	}
	chipsVal, ok := chipsI.([]interface{})
	if !ok {
		tools.Logger.Fatal("chips assert failure")
	}
	for _, v := range chipsVal {
		chipF, ok := v.(float64)
		if !ok {
			tools.Logger.Fatal("chipF assert failure")
		}
		chips = append(chips, int32(chipF))
	}
	tools.Logger.Infof("chips: %v", chips)
}

func parseBetZones(res *models.Response) {
	betZonesI, ok := res.Data["bet_zones"]
	if !ok {
		tools.Logger.Fatal("bet zones is required")
	}
	betZonesVal, ok := betZonesI.([]interface{})
	if !ok {
		tools.Logger.Fatal("betZonesVal assert failure")
	}
	for _, v := range betZonesVal {
		betZoneF, ok := v.(float64)
		if !ok {
			tools.Logger.Fatal("betZoneF assert failure")
		}
		betZones = append(betZones, int32(betZoneF))
	}
	tools.Logger.Infof("betZones: %v", betZones)
}
