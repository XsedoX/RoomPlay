package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance"
	"github.com/XsedoX/RoomPlay/infrastructure/validation"
	"github.com/XsedoX/RoomPlay/initialization"
	"github.com/XsedoX/RoomPlay/presentation"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
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
	validation.Initialize()
	ctx := context.Background()

	configuration := config.Load()

	db := persistance.InitializeDatabase(ctx, configuration)
	if configuration.IsDevelopment() {
		tableExists, _ := tableExists(db, "users")
		if tableExists {
			log.Print("Schema already applied skipping.")
		} else { // Seed database once
			projectRoot, err := integration_tests.FindProjectRoot()
			if err != nil {
				log.Fatalf("failed to find project root: %v", err)
			}

			schemaPath := filepath.Join(projectRoot, "infrastructure", "persistance", "RoomPlay2.sql")
			if err := integration_tests.ApplySchema(ctx, schemaPath, db); err != nil {
				log.Fatalf("failed to apply schema: %v", err)
			}
			log.Print("Applied database schema.")
		}
	}
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

func tableExists(db *sqlx.DB, tableName string) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' AND table_name = $1
        );
    `
	err := db.QueryRow(query, tableName).Scan(&exists)
	return exists, err
}
