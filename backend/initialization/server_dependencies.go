package initialization

import (
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/application/room"
	"xsedox.com/main/application/user"
	"xsedox.com/main/infrastructure/persistance"
	"xsedox.com/main/presentation/handlers"
)

type ServerDependencies struct {
	roomHandler    *handlers.RoomHandler
	userHandler    *handlers.UserHandler
	oidcHandler    *handlers.OidcHandler
	userRepository user.IRepository
	roomRepository room.IRepository
}

func NewServerDependencies(db *sqlx.DB) *ServerDependencies {
	userRepository := persistance.NewUserRepository(db)
	roomRepository := persistance.NewRoomRepository(db)
	unitOfWork := persistance.NewUnitOfWork(db)

	loginUserCommandHandler := user.NewLoginCommandHandler(userRepository, unitOfWork)
	userHandler := handlers.NewUserHandler(loginUserCommandHandler)
	oidcHandler := handlers.NewOidcHandler(loginUserCommandHandler)

	createRoomCommandHandler := room.NewCreateCommandHandler(roomRepository, unitOfWork)
	roomHandler := handlers.NewRoomHandler(createRoomCommandHandler)

	return &ServerDependencies{
		userHandler:    userHandler,
		roomHandler:    roomHandler,
		oidcHandler:    oidcHandler,
		userRepository: userRepository,
		roomRepository: roomRepository,
	}
}
func (sd ServerDependencies) RoomHandler() *handlers.RoomHandler {
	return sd.roomHandler
}
func (sd ServerDependencies) UserHandler() *handlers.UserHandler {
	return sd.userHandler
}
func (sd ServerDependencies) OidcHandler() *handlers.OidcHandler {
	return sd.oidcHandler
}
