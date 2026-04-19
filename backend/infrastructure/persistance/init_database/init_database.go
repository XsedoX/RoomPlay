package init_database

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"path/filepath"
	"sort"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//go:embed sql_scripts/*.sql
var sqlScripts embed.FS

func InitializeDatabase(ctx context.Context, connectionString string) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, "pgx", connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	entries, err := sqlScripts.ReadDir("sql_scripts")
	if err != nil {
		log.Fatalf("Unable to read sql_scripts directory: %v\n", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	txx, err := db.Beginx()
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v\n", err)
	}
	handleRollback := func() {
		if err := txx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Fatalf("Transaction rollback failed: %v\n", err)
		}
	}
	defer handleRollback()

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		content, err := sqlScripts.ReadFile(filepath.Join("sql_scripts", entry.Name()))
		if err != nil {
			log.Fatalf("Unable to read SQL script %s: %v\n", entry.Name(), err)
		}

		_, err = txx.ExecContext(ctx, string(content))
		if err != nil {
			log.Fatalf("Unable to execute SQL script %s: %v\n", entry.Name(), err)
		}

		log.Printf("Executed SQL script: %s\n", entry.Name())
	}

	// seeder := seeder.NewSeeder(txx)
	// if err := seeder.SeedAll(ctx); err != nil {
	// 	log.Fatalf("failed to seed database: %v", err)
	// }

	if err := txx.Commit(); err != nil {
		log.Fatalf("Transaction commit failed: %v\n", err)
	}

	return db
}
