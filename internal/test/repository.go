package test

import "context"

type Repository interface {
	IncrementByValue(ctx context.Context, key string, value int64) (int64, error)
}
