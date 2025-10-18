package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"xsedox.com/application/room"
	"xsedox.com/infrastructure/persistance"
	_ "xsedox.com/main/docs"
	"xsedox.com/presentation"

	"log"
)

// @title RoomPlay API
// @version 1.0
// @description This is the API for the RoomPlay service.
// @host localhost:7865
// @BasePath /
func main() {
	ctx := context.Background()
	connString := "postgres://postgres:dupa1234@localhost:5432/postgres"
	db, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer func(db *pgx.Conn, ctx context.Context) {
		err := db.Close(ctx)
		if err != nil {
			log.Fatalf("Unable to close database connection: %v\n", err)
		}
	}(db, context.Background())
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	// Dependencies
	roomRepository := persistance.NewRoomRepository(db, &ctx)
	createRoomHandler := room.NewCreateCommandHandler(roomRepository)
	roomHandler := presentation.NewRoomHandler(createRoomHandler)

	// Start Server
	server := presentation.NewServer(roomHandler)
	server.Start()
}
