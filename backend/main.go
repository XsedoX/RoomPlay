package main

import (
	"context"
	"log"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/init_database"
	"github.com/XsedoX/RoomPlay/initialization/initialize_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/api_server"
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
	ctx := context.Background()

	configuration := config.Load()

	db := init_database.InitializeDatabase(ctx, configuration.Database().ConnectionString)

	dependencies := initialize_dependencies.NewServerDependencies(db, configuration)

	log.Printf("Loaded config: port: %v, host: %v, environment: %v", configuration.Server().Port, configuration.Server().Host, configuration.Environment)

	defer func(db *sqlx.DB, ctx context.Context) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database connection: %v\n", err)
		}
	}(db, context.Background())
	// Start Server
	server := api_server.NewServer(dependencies)
	server.Start(configuration)
}
