package contracts

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type IQueryer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)               // insert / update / delete
	NamedExecContext(ctx context.Context, query string, args interface{}) (sql.Result, error)     // insert / update / delete batched
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error    // single row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error // multiple rows
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)             // manual
}
