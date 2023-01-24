package repositories

import (
	tools "github.com/duel80003/my-tools"
	"github.com/gomodule/redigo/redis"
	"main-service/drivers"
)

func GetPlayerGameId(sid string) string {
	gid, err := redis.String(drivers.GetRedisConn().Do("HGET", sid, "game_id"))
	if err != nil {
		if redis.ErrNil == err {
			return ""
		}
		tools.Logger.Errorf("[GetPlayerGameId] error: %s", err)
		return ""
	}
	tools.Logger.Debugf("[GetPlayerGameId] gid: %s", gid)
	return gid
}
