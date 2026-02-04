package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/init_database"
	"github.com/XsedoX/RoomPlay/initialization/initialize_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/api_server"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
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

	db := init_database.InitializeDatabase(ctx, configuration)
	if configuration.IsDevelopment() {
		tableExists, _ := tableExists(db, "users")
		if tableExists {
			log.Print("Schema already applied skipping.")
		} else { // Seed database once
			projectRoot, err := tests_initializer.FindProjectRoot()
			if err != nil {
				log.Fatalf("failed to find project root: %v", err)
			}

			schemaPath := filepath.Join(projectRoot, "infrastructure", "persistance", "RoomPlay2.sql")
			if err := tests_initializer.ApplySchema(ctx, schemaPath, db); err != nil {
				log.Fatalf("failed to apply schema: %v", err)
			}
			log.Print("Applied database schema.")
		}
	}
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
