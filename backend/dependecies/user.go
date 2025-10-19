package dependecies

import (
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/application/user"
	"xsedox.com/main/infrastructure/persistance"
	"xsedox.com/main/presentation/handlers"
)

type UserDependencies struct {
	userRepository    *persistance.UserRepository
	createUserHandler *user.CreateCommandHandler
	userHandler       *handlers.UserHandler
	unitOfWork        *persistance.UnitOfWork
}

func NewUserDependencies(db *sqlx.DB) *UserDependencies {
	userRepository := persistance.NewUserRepository(db)
	unitOfWork := persistance.NewUnitOfWork(db)
	createUserHandler := user.NewCreateCommandHandler(userRepository, unitOfWork)
	userHandler := handlers.NewUserHandler(createUserHandler)

	return &UserDependencies{
		userRepository:    userRepository,
		createUserHandler: createUserHandler,
		userHandler:       userHandler,
		unitOfWork:        unitOfWork,
	}
}
func (userDeps *UserDependencies) GetUserHandler() *handlers.UserHandler {
	return userDeps.userHandler
}
