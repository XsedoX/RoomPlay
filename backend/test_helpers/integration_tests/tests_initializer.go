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
	Ctx          context.Context
	InjectedUser = SeedData.Users[0]
)

func InitializeDatabaseContainer() {
	Ctx = context.Background()
	var err error
	PgContainer, err = SetupPostgres(Ctx)
	if err != nil {
		log.Fatalf("failed to setup postgres: %v", err)
	}

	projectRoot, err := FindProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project root: %v", err)
	}

	schemaPath := filepath.Join(projectRoot, "infrastructure", "persistance", "RoomPlay2.sql")
	if err := ApplySchema(Ctx, schemaPath, PgContainer.DB); err != nil {
		log.Fatalf("failed to apply schema: %v", err)
	}

	// Seed database once
	dbx := PgContainer.DB
	seeder := NewSeeder(dbx)
	if err := seeder.SeedAll(Ctx); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}
}

func RunTestsWithDatabase(m *testing.M) {
	code := m.Run()
	if err := PgContainer.Teardown(Ctx); err != nil {
		log.Printf("failed to teardown postgres: %v", err)
	}

	os.Exit(code)
}

func GetTxxAndCtx(t *testing.T) (*sqlx.Tx, context.Context) {
	t.Helper()
	require.NotNil(t, PgContainer, "pgContainer is nil; TestMain likely didn’t run")
	ctx := context.Background()
	dbx := PgContainer.DB
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = txx.Rollback() })
	return txx, Ctx
}

func InitializeApiServer(m *testing.M) {
	validation.Initialize()

	InitializeDatabaseContainer()
	configuration := othermocks.MockConfiguration{}

	db := PgContainer.DB
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
