package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"xsedox.com/main/config"
	"xsedox.com/main/infrastructure/persistance"
	"xsedox.com/main/infrastructure/validation"
	"xsedox.com/main/initialization"
	"xsedox.com/main/presentation"
)

import (
	"context"

	_ "xsedox.com/main/docs"
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
