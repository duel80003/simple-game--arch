package drivers

import (
	tools "github.com/duel80003/my-tools"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

var (
	pool *redis.Pool
)

func RedisInit() {
	host := os.Getenv("REDIS_HOST")
	pool = &redis.Pool{
		MaxIdle:     20,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host)
		},
	}
	tools.Logger.Infof("redis connected")
}

func RedisFlushAll() {
	_, err := pool.Get().Do("FLUSHALL")
	if err != nil {
		tools.Logger.Errorf("RedisFlushAll error: %s", err)
	}
}

func RedisClose() {
	err := pool.Close()
	if err != nil {
		tools.Logger.Errorf("redis disconnect faulure, error: %s", err)
	}
	tools.Logger.Infof("redis disconnected")
}

func GetRedisConn() redis.Conn {
	return pool.Get()
}
