package repositories

import (
	"game-process-service/config"
	"game-process-service/drivers"
	tools "github.com/duel80003/my-tools"
	"github.com/joho/godotenv"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("../.env")
	tools.LogInit()
	if err != nil {
		tools.Logger.Infof("load env file failure")
	}
	drivers.RedisInit()

}

func TestUpdateRoomBetInfo(t *testing.T) {
	defer drivers.RedisClose()
	rid := "test"
	var betZone1 int32 = 0
	var betZone2 int32 = 1
	var bet int32 = 5
	UpdateRoomBetInfo(rid, config.BetZoneMap[betZone1], bet)
	UpdateRoomBetInfo(rid, config.BetZoneMap[betZone2], bet)
	time.Sleep(1 * time.Second)
}

func TestResetRoomBetInfo(t *testing.T) {
	defer drivers.RedisClose()
	rid := "test"
	ResetRoomBetInfo(rid, []string{config.BetZoneMap[0], config.BetZoneMap[1]})
	time.Sleep(1 * time.Second)
}

func TestGetRoomBetInfo(t *testing.T) {
	defer drivers.RedisClose()
	rid := "test"
	betZones := GetRoomBetInfo(rid)
	t.Logf("bet zones: %+v \n", betZones)
	time.Sleep(1 * time.Second)
}
