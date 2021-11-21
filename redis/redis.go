package redis

import "github.com/go-redis/redis/v8"

var (
	redisCli *redis.Client
)

func NewRedisClient() {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
	})
}

func GetRedisCli() *redis.Client {
	return redisCli
}
