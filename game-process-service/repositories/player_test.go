package repositories

import (
	"game-process-service/config"
	"game-process-service/drivers"
	"github.com/google/uuid"
	"testing"
	"time"
)

var (
	sid string = "b1126668-b9eb-45b6-b82a-e6d063932be4"
)

func TestSetPlayer(t *testing.T) {
	defer drivers.RedisClose()
	sid = uuid.NewString()
	pid := "test"
	SetPlayer(sid, pid, "01")
	time.Sleep(1 * time.Second)
}

func TestUpdatePlayerBetInfo(t *testing.T) {
	defer drivers.RedisClose()
	UpdatePlayerBetInfo(sid, config.BetZoneMap[1], 20)
	time.Sleep(1 * time.Second)
}

func TestGetPlayer(t *testing.T) {
	defer drivers.RedisClose()
	p := GetPlayer(sid)
	t.Logf("player: %+v", p)
}

func TestResetPlayerBetInfo(t *testing.T) {
	defer drivers.RedisClose()
	ResetPlayerBetInfo([]string{sid}, []string{config.BetZoneMap[0], config.BetZoneMap[1]})
	time.Sleep(1 * time.Second)
}

func TestRemovePlayer(t *testing.T) {
	defer drivers.RedisClose()
	RemovePlayer(sid)
	time.Sleep(1 * time.Second)
}
