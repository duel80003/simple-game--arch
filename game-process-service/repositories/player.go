package repositories

import (
	"game-process-service/drivers"
	"game-process-service/models"
	tools "github.com/duel80003/my-tools"
	"github.com/gomodule/redigo/redis"
)

func SetPlayer(sid, pid, gid string) {
	conn := drivers.GetRedisConn()
	conn.Send("HSET", sid, "player_id", pid)
	conn.Send("HSET", sid, "game_id", gid)
	conn.Flush()
}

func UpdatePlayerBetInfo(sid, betZone string, bet int32) {
	tools.Logger.Debugf("[UpdatePlayerBetInfo] sid: %s, betZone: %s, bet: %d", sid, betZone, bet)
	do, err := drivers.GetRedisConn().Do("HINCRBY", sid, betZone, bet)
	if err != nil {
		tools.Logger.Errorf("[UpdatePlayerBetInfo] error: %s", err)
		return
	}
	tools.Logger.Debugf("[UpdatePlayerBetInfo] do: %+v", do)
}

func ResetPlayerBetInfo(sids, betZones []string) {
	conn := drivers.GetRedisConn()
	for _, sid := range sids {
		for _, v := range betZones {
			conn.Send("HSET", sid, v, 0)
		}
	}
	conn.Flush()
}

func GetPlayer(sid string) (p *models.Player) {
	values, err := redis.Values(drivers.GetRedisConn().Do("HGETALL", sid))
	if err != nil {
		tools.Logger.Errorf("[GetRoomBetInfo] get info error: %s", err)
		return
	}
	p = new(models.Player)
	err = redis.ScanStruct(values, p)
	return
}

func RemovePlayer(sid string) {
	drivers.GetRedisConn().Do("DEL", sid)
}
