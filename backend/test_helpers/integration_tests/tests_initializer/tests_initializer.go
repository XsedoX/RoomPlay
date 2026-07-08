package tests_initializer

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/init_database"
	"github.com/XsedoX/RoomPlay/presentation/api_server"
	"github.com/XsedoX/RoomPlay/presentation/application_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/initialize_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/presentation_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_configuration"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	PgContainer  *PostgresContainer
	TestServer   *api_server.Server
	ctx          context.Context
	InjectedUser = seeder.SeedData.Users[0]
)

func InitializeDatabaseContainer() {
	ctx = context.Background()
	var err error
	PgContainer, err = SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to setup postgres: %v", err)
	}

	init_database.InitializeDatabase(ctx, PgContainer.connStr)
	// Seed database once
	dbx := PgContainer.db
	txx, err := dbx.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatalf("failed to begin transaction for seeding database: %v", err)
	}
	defer txx.Rollback()
	seeder := seeder.NewSeeder(txx)
	if err := seeder.SeedAll(ctx); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}
	txx.Commit()
}

func RunTestsWithDatabase(m *testing.M) {
	code := m.Run()
	if err := PgContainer.Teardown(ctx); err != nil {
		log.Printf("failed to teardown postgres: %v", err)
	}
	os.Exit(code)
}

func GetTxxAndCtx(t *testing.T, reinitDb bool) (*sqlx.Tx, context.Context) {
	t.Helper()
	require.NotNil(t, PgContainer, "pgContainer is nil; TestMain likely didn’t run")
	dbx := PgContainer.db
	txx, err := dbx.BeginTxx(ctx, nil)
	require.NoError(t, err)
	if !reinitDb {
		t.Cleanup(func() { _ = txx.Rollback() })
	} else {
		t.Cleanup(func() {
			_ = txx.Rollback()
			reseedDatabase(t)
		})
	}
	return txx, ctx
}

func reseedDatabase(t *testing.T) {
	db := PgContainer.db
	_, err := db.Exec(`DO $$ 
DECLARE 
    r RECORD;
BEGIN
    -- Disable all triggers (including foreign key constraints)
    SET CONSTRAINTS ALL DEFERRED;
    
    FOR r IN (SELECT tablename, schemaname FROM pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema')) LOOP
        EXECUTE 'ALTER TABLE ' || quote_ident(r.schemaname) || '.' || quote_ident(r.tablename) || ' DISABLE TRIGGER ALL';
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.schemaname) || '.' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
        EXECUTE 'ALTER TABLE ' || quote_ident(r.schemaname) || '.' || quote_ident(r.tablename) || ' ENABLE TRIGGER ALL';
    END LOOP;
END $$;`)
	require.NoError(t, err)
	// Seed database once
	txx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatalf("failed to begin transaction for reseeding database: %v", err)
	}
	defer txx.Rollback()
	seeder := seeder.NewSeeder(txx)
	if err := seeder.SeedAll(ctx); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}
	txx.Commit()
}

func InitializeApiServer(m *testing.M) {
	setup_validation.Initialize()

	InitializeDatabaseContainer()
	configuration := mock_configuration.MockConfiguration{}
	mockMusicDataService := mock_music_data_provider_service.MockMusicDataProviderService{}
	mockMusicDataService.On(
		"SearchSongsByQuery",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*string"),
		mock.AnythingOfType("uint8"),
	).Return(seeder.ExternalSongData, nil)

	db := PgContainer.db
	InjectedUserId := InjectedUser.Id()
	injectedUserClaim := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), user.IdClaimContextKeyName, &InjectedUserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	infrastructureDependencies := infrastructure_dependencies.ConstructInfrastructureDependencies(
		ctx,
		db,
		&configuration,
	)
	infrastructureDependencies.CachingSongDecorator = &mockMusicDataService
	applicationDependencies := application_dependencies.ConstructApplicationDependencies(
		infrastructureDependencies,
		&configuration,
	)
	presentationDependencies := presentation_dependencies.ConstructPresentationDependencies(
		&configuration,
		applicationDependencies,
		infrastructureDependencies,
	)
	dependencies := initialize_dependencies.NewServerDependencies(
		infrastructureDependencies,
		presentationDependencies,
	)

	server := api_server.NewServer(dependencies, &configuration, injectedUserClaim)
	TestServer = server
	RunTestsWithDatabase(m)
}
