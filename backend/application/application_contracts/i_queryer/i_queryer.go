package i_queryer

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type IQueryer interface {
	// NOTE: Use for INSERT, UPDATE, or DELETE statements where you don't need to scan any rows back.
	// Returns sql.Result which gives you LastInsertId() and RowsAffected(). Args are positional ($1, $2, ...).
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// NOTE: Use for INSERT, UPDATE, or DELETE with named parameters (e.g. :field_name) bound from a struct or map.
	// Particularly useful for batch inserts where you pass a slice of structs.
	NamedExecContext(ctx context.Context, query string, args interface{}) (sql.Result, error)

	// NOTE: Use when you expect exactly ONE row back (e.g. SELECT ... WHERE id = $1).
	// Scans the result directly into dest (a struct pointer). Returns sql.ErrNoRows if no row is found.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// NOTE: Use when you expect MULTIPLE rows back (e.g. SELECT ... WHERE active = true).
	// Scans all results into dest (a pointer to a slice of structs). Returns an empty slice if no rows match.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// NOTE: Use when you need manual control over row iteration (e.g. streaming large result sets or
	// complex per-row logic). You must call rows.Next(), rows.StructScan()/rows.Scan(), and rows.Close() yourself.
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)

	// NOTE: Use when you need a single raw *sql.Row for manual scanning (e.g. scanning into non-struct types
	// or when you need fine-grained control). You must call row.Scan() yourself.
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
}
