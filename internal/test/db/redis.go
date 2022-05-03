package db

import (
	"console-application-service-age/internal/test"
	"console-application-service-age/pkg/client/redis"
	"context"
)

type repository struct {
	client redis.Client
}

func (r repository) IncrementByValue(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

func NewRepository(client redis.Client) test.Repository {
	return &repository{client: client}
}
