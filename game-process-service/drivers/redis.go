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
			conn, err := redis.Dial("tcp", host)
			if err != nil {
				tools.Logger.Errorf("redis dial error: %s", err)
				return nil, err
			}
			return conn, err
		},
	}
	tools.Logger.Infof("redis connected")
}

func RedisFlushDB() {
	_, err := pool.Get().Do("FLUSHDB")
	if err != nil {
		tools.Logger.Errorf("RedisFlushDB error: %s", err)
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
