package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"xsedox.com/main/config"
	"xsedox.com/main/infrastructure/oidc"
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
	err := validation.Initialize()
	if err != nil {
		log.Fatalf("failed to register validation: %v", err)
	}

	config.Config = config.Load()
	log.Printf("Loaded config: port: %v, host: %v, environment: %v", config.Config.Server.Port, config.Config.Server.Host, config.Config.Environment)
	ctx := context.Background()
	db := persistance.InitializeDatabase(ctx)

	defer func(db *sqlx.DB, ctx context.Context) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database connection: %v\n", err)
		}
	}(db, context.Background())
	// Dependencies
	serverDependencies := initialization.NewServerDependencies(db)

	fmt.Printf("This is the hash: %x\n", *oidc.NewEncryptionKey())
	fmt.Printf("This is the hash: %x\n", *oidc.NewEncryptionKey())
	fmt.Printf("This is the hash: %x\n", *oidc.NewEncryptionKey())
	// Start Server
	server := presentation.NewServer(serverDependencies)
	server.Start()
}
