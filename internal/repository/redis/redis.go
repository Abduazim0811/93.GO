package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func SetupRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return rdb
}