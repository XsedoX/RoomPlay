package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/init_database"
	"github.com/XsedoX/RoomPlay/presentation/api_server"
	"github.com/XsedoX/RoomPlay/presentation/application_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/initialize_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/presentation_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	_ "github.com/XsedoX/RoomPlay/docs"
)

// @title RoomPlay API
// @version 1.0
// @description This is the API for the RoomPlay service.
// @host localhost:7654
// @BasePath /api/v1
//
// Security definition for Bearer JWT in the Authorization header
// @securityDefinitions.apikey BearerAuth
// @in cookie
// @name Authorization
// @description Type "Bearer {token}"
func main() {
	setup_validation.Initialize()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configuration := config.Load()
	db := init_database.InitializeDatabase(ctx, configuration.Database().ConnectionString)

	redisConnStr := strings.Replace(
		configuration.Database().RedisConnectionString,
		"valkey", "redis", 1,
	)
	opts, err := redis.ParseURL(redisConnStr)
	if err != nil {
		log.Fatalf("Unable to parse Redis connection string: %v\n", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
	}

	infrastructureDependendies := infrastructure_dependencies.ConstructInfrastructureDependencies(
		ctx,
		db,
		rdb,
		configuration,
	)
	applicationDependencies := application_dependencies.ConstructApplicationDependencies(
		infrastructureDependendies,
		configuration,
	)
	presentationDependencies := presentation_dependencies.ConstructPresentationDependencies(
		configuration,
		applicationDependencies,
		infrastructureDependendies,
	)
	dependencies := initialize_dependencies.NewServerDependencies(
		infrastructureDependendies,
		presentationDependencies,
	)

	log.Printf("Loaded config: port: %v, host: %v, environment: %v", configuration.Server().Port, configuration.Server().Host, configuration.Environment)

	defer func(db *sqlx.DB, ctx context.Context) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database connection: %v\n", err)
		}
	}(db, context.Background())
	// Start Server
	server := api_server.NewServer(dependencies, configuration)
	server.Start(configuration)
}
