package tests_initializer

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	container *postgres.PostgresContainer
	db        *sqlx.DB
	connStr   string
}

func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	dbName := "roomplay"
	dbUser := "user"
	dbPassword := "password"

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(40*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &PostgresContainer{
		container: pgContainer,
		db:        db,
		connStr:   connStr,
	}, nil
}

func (pc *PostgresContainer) Teardown(ctx context.Context) error {
	if pc.db != nil {
		pc.db.Close()
	}
	if pc.container != nil {
		return pc.container.Terminate(ctx)
	}
	return nil
}
