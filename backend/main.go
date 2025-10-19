package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"xsedox.com/main/config"
	"xsedox.com/main/dependecies"
	"xsedox.com/main/presentation"
	"xsedox.com/main/validation"
)

import (
	"context"

	_ "github.com/jackc/pgx/stdlib"
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

	conf := config.Load()
	log.Printf("Loaded config: port: %v, host: %v, environment: %v", conf.Server.Port, conf.Server.Host, conf.Environment)
	ctx := context.Background()
	db, err := sqlx.ConnectContext(ctx, "pgx", conf.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer func(db *sqlx.DB, ctx context.Context) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database connection: %v\n", err)
		}
	}(db, context.Background())

	// Dependencies
	serverDependencies := &presentation.ServerDependencies{
		UserHandler:   dependecies.NewUserDependencies(db).GetUserHandler(),
		RoomHandler:   dependecies.NewRoomDependencies(db).GetRoomHandler(),
		Configuration: conf,
	}
	// Start Server
	server := presentation.NewServer(serverDependencies)
	server.Start()
}
