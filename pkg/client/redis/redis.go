package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Client interface {
	IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd
}

func NewClient(host, port string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})
}
