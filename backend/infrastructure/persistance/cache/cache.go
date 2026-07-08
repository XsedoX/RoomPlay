package cache

import (
	"context"
	"encoding/json"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
)

type ICache[T any] interface {
	Get(key string, ctx context.Context, queryer i_queryer.IQueryer) (T, error)
	Set(key string, value T, ctx context.Context, queryer i_queryer.IQueryer) error
	Remove(key string, ctx context.Context, queryer i_queryer.IQueryer) error
}
type Cache[T any] struct {
	cacheSimilarityThreshold float32
}

func NewCache[T any](cacheSimilarityThreshold float32) *Cache[T] {
	return &Cache[T]{cacheSimilarityThreshold: cacheSimilarityThreshold}
}

func (c *Cache[T]) Get(key string, ctx context.Context, queryer i_queryer.IQueryer) (T, error) {
	var value json.RawMessage
	err := queryer.GetContext(ctx, &value, `
		SELECT value FROM cache
		WHERE similarity(key, $1) > $2
		ORDER BY similarity(key, $1) DESC
		LIMIT 1;
		`,
		key,
		c.cacheSimilarityThreshold,
	)
	var zero T
	if err != nil {
		return zero, err
	}
	var result T
	err = json.Unmarshal(value, &result)
	if err != nil {
		return zero, err
	}
	return result, nil
}

func (c *Cache[T]) Set(key string, value T, ctx context.Context, queryer i_queryer.IQueryer) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = queryer.ExecContext(ctx, `
		INSERT INTO cache(key, value)
		VALUES($1, $2::jsonb)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value
	`, key, b)
	return err
}

func (c *Cache[T]) Remove(key string, ctx context.Context, queryer i_queryer.IQueryer) error {
	_, err := queryer.ExecContext(ctx, `DELETE FROM cache WHERE key = $1`, key)
	return err
}
