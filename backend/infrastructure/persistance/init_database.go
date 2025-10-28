package persistance

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/config"
)

func InitializeDatabase(ctx context.Context, configuration config.IConfiguration) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, "pgx", configuration.Database().ConnectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	_, err = db.ExecContext(ctx, `
		CREATE UNLOGGED TABLE IF NOT EXISTS cache(
		id SERIAL PRIMARY KEY,
		key TEXT UNIQUE NOT NULL,
		value JSONB,
		created_at_utc TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
		
		CREATE INDEX IF NOT EXISTS cache_idx ON cache(id) INCLUDE (value);`)
	if err != nil {
		log.Fatalf("Unable to create cache table: %v\n", err)
	}
	return db
}
