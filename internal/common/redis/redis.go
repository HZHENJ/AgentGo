package redis

import (
	"agentgo/pkg/conf"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() error {
	c := conf.Config.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPw,
		DB:       c.RedisDb,
	})

	// Best Practice: Create a context with timeout for the Ping operation at startup (e.g., 5 seconds)
	// This way, if the Redis server is unreachable, the program won't hang indefinitely during startup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() 

	_, err := RDB.Ping(ctx).Result()
	return err
}