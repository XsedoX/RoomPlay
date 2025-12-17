package persistance

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"xsedox.com/main/test_helpers/infrastructure_test"
)

var (
	pgContainer *infrastructure_test.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgContainer, err = infrastructure_test.SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to setup postgres: %v", err)
	}

	projectRoot, err := infrastructure_test.FindProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	schemaPath := filepath.Join(projectRoot, "infrastructure", "persistance", "RoomPlay2.sql")
	if err := pgContainer.ApplySchema(ctx, schemaPath); err != nil {
		log.Fatalf("failed to apply schema: %v", err)
	}

	// Seed database once
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	seeder := infrastructure_test.NewSeeder(dbx)
	if err := seeder.SeedAll(ctx); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	code := m.Run()

	if err := pgContainer.Teardown(ctx); err != nil {
		log.Printf("failed to teardown postgres: %v", err)
	}

	os.Exit(code)
}

func GetTxxAndCtx(t *testing.T) (*sqlx.Tx, context.Context) {
	t.Helper()
	require.NotNil(t, pgContainer, "pgContainer is nil; TestMain likely didn’t run")
	ctx := context.Background()
	dbx := sqlx.NewDb(pgContainer.DB, "pgx")
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = txx.Rollback() })
	return txx, ctx
}
