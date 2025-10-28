package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"xsedox.com/main/config"
	"xsedox.com/main/infrastructure/persistance"
	"xsedox.com/main/initialization"
	"xsedox.com/main/presentation"
	"xsedox.com/main/validation"
)

import (
	"context"

	_ "xsedox.com/main/docs"
)

// @title RoomPlay API
// @version 1.0
// @description This is the API for the RoomPlay service.
// @host localhost:7865
// @BasePath /
//
// Security definition for Bearer JWT in the Authorization header
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {token}"
func main() {
	ctx := context.Background()
	err := validation.Initialize()
	configuration := config.Load()

	db := persistance.InitializeDatabase(ctx, configuration)

	dependencies := initialization.NewServerDependencies(db, configuration)

	if err != nil {
		log.Fatalf("failed to register validation: %v", err)
	}

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
