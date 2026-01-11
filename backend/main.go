package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance"
	"github.com/XsedoX/RoomPlay/infrastructure/validation"
	"github.com/XsedoX/RoomPlay/initialization"
	"github.com/XsedoX/RoomPlay/presentation"
)

import (
	"context"

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
	validation.Initialize()
	ctx := context.Background()

	configuration := config.Load()

	db := persistance.InitializeDatabase(ctx, configuration)
	
	dependencies := initialization.NewServerDependencies(db, configuration)

	log.Printf("Loaded config: port: %v, host: %v, environment: %v", configuration.Server().Port, configuration.Server().Host, configuration.Environment)

	defer func(db *sqlx.DB, ctx context.Context) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database connection: %v\n", err)
		}
	}(db, context.Background())
	// Start Server
	server := presentation.NewServer(dependencies)
	server.Start(configuration)
}
