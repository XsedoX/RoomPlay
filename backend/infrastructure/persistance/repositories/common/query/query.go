package query

import (
	"context"

	"gorm.io/gorm"
)

type Scope = func(*gorm.DB) *gorm.DB

func Exists[T any](ctx context.Context, db *gorm.DB, scopes ...Scope) (bool, error) {
	count, err := gorm.G[T](db.Scopes(scopes...)).
		Select("1").
		Limit(1).
		Count(
			ctx,
			"1",
		)

	return count > 0, err
}
