package persistance

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type ICache[T any] interface {
	Get(key string, ctx context.Context) (T, error)
	Set(key string, value T, ctx context.Context) error
	Remove(key string, ctx context.Context) error
}
type Cache[T any] struct {
	db *sqlx.DB
}

func (c *Cache[T]) Get(key string, ctx context.Context) (*T, error) {
	var value json.RawMessage
	err := c.db.GetContext(ctx, &value, `SELECT value FROM cache WHERE key = $1;`,
		key)
	if err != nil {
		return nil, err
	}
	var result T
	err = json.Unmarshal(value, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Cache[T]) Set(key string, value T, ctx context.Context) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = c.db.ExecContext(ctx, `
		INSERT INTO cache(key, value)
		VALUES($1, $2::jsonb)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value
	`, key, b)
	return err
}

func (c *Cache[T]) Remove(key string, ctx context.Context) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM cache WHERE key = $1`, key)
	return err
}
