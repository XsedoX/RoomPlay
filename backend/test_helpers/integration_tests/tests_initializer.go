package integration_tests

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/infrastructure/validation"
	"github.com/XsedoX/RoomPlay/initialization"
	"github.com/XsedoX/RoomPlay/presentation"
	othermocks "github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	PgContainer  *PostgresContainer
	TestServer   *presentation.Server
	ctx          context.Context
	InjectedUser = SeedData.Users[0]
)

func InitializeDatabaseContainer() {
	ctx = context.Background()
	var err error
	PgContainer, err = SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to setup postgres: %v", err)
	}

	projectRoot, err := FindProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	schemaPath := filepath.Join(projectRoot, "infrastructure", "persistance", "RoomPlay2.sql")
	if err := ApplySchema(ctx, schemaPath, PgContainer.db); err != nil {
		log.Fatalf("failed to apply schema: %v", err)
	}

	// Seed database once
	dbx := PgContainer.db
	seeder := NewSeeder(dbx)
	if err := seeder.SeedAll(ctx); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}
}

func RunTestsWithDatabase(m *testing.M) {
	code := m.Run()
	if err := PgContainer.Teardown(ctx); err != nil {
		log.Printf("failed to teardown postgres: %v", err)
	}
	os.Exit(code)
}

func GetTxxAndCtx(t *testing.T) (*sqlx.Tx, context.Context) {
	t.Helper()
	require.NotNil(t, PgContainer, "pgContainer is nil; TestMain likely didn’t run")
	dbx := PgContainer.db
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = txx.Rollback() })
	return txx, ctx
}

func InitializeApiServer(m *testing.M) {
	validation.Initialize()

	InitializeDatabaseContainer()
	configuration := othermocks.MockConfiguration{}

	db := PgContainer.db
	InjectedUserId := InjectedUser.Id()
	injectedUserClaim := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), user.IdClaimContextKeyName, &InjectedUserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	dependencies := initialization.NewServerDependencies(db, &configuration)
	server := presentation.NewServer(dependencies, injectedUserClaim)
	TestServer = server
	RunTestsWithDatabase(m)
}
