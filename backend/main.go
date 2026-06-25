package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/init_database"
	"github.com/XsedoX/RoomPlay/presentation/api_server"
	"github.com/XsedoX/RoomPlay/presentation/initialize_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/jmoiron/sqlx"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db := init_database.InitializeDatabase(ctx, config.Load().Database().ConnectionString)
	configuration := config.Load()

	dependencies := initialize_dependencies.NewServerDependencies(ctx, db, configuration)

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
