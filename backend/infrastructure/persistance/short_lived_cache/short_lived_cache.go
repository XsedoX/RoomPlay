package short_lived_cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type IShortLivedCache[T any] interface {
	Set(key string, value T, ctx context.Context) error
	Get(key string, ctx context.Context) (T, error)
}

type ShortLivedCache[T any] struct {
	rdb *redis.Client
}

func NewShortLivedCache[T any](rdb *redis.Client) *ShortLivedCache[T] {
	return &ShortLivedCache[T]{
		rdb: rdb,
	}
}

func (c *ShortLivedCache[T]) Set(key string, value T, ctx context.Context) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, 20*time.Second).Err()
}

func (c *ShortLivedCache[T]) Get(key string, ctx context.Context) (T, error) {
	var value T
	result, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return value, err
	}

	if err := json.Unmarshal([]byte(result), &value); err != nil {
		return value, err
	}

	return value, nil
}
