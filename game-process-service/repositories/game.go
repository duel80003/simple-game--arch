package repositories

import (
	"game-process-service/drivers"
	"game-process-service/models"
	tools "github.com/duel80003/my-tools"
	"github.com/gomodule/redigo/redis"
)

func UpdateRoomBetInfo(rid, betZone string, bet int32) {
	tools.Logger.Debugf("[UpdateRoomBetInfo] rid: %s, betZone: %s, bet: %d", rid, betZone, bet)
    conn := drivers.GetRedisConn()
    defer conn.Close()
    do, err := conn.Do("HINCRBY", rid, betZone, bet)
	if err != nil {
		tools.Logger.Errorf("[UpdateRoomBetInfo] error: %s", err)
		return
	}
	tools.Logger.Debugf("[UpdateRoomBetInfo] do: %+v", do)
}

func ResetRoomBetInfo(rid string, betZones []string) {
	conn := drivers.GetRedisConn()
    defer conn.Close()
	for _, v := range betZones {
		conn.Send("HSET", rid, v, 0)
	}
	conn.Flush()
}

func GetRoomBetInfo(rid string) (betZones *models.BetZones) {
	conn := drivers.GetRedisConn()
    defer conn.Close()
	values, err := redis.Values(conn.Do("HGETALL", rid))
	if err != nil {
		tools.Logger.Errorf("[GetRoomBetInfo] get info error: %s", err)
		return
	}

	betZones = new(models.BetZones)
	err = redis.ScanStruct(values, betZones)
	if err != nil {
		tools.Logger.Errorf("ScanStruct error: %s", err)
		return
	}
	tools.Logger.Debugf("result: %+v", betZones)
	return
}
