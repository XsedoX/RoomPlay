package initialization

import (
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/join"
	"xsedox.com/main/application/services"
	"xsedox.com/main/application/user/data"
	"xsedox.com/main/application/user/login"
	"xsedox.com/main/application/user/register"
	"xsedox.com/main/config"
	"xsedox.com/main/infrastructure/authentication"
	"xsedox.com/main/infrastructure/persistance"
	"xsedox.com/main/presentation/controllers"
)

type ServerDependencies struct {
	roomController *controllers.RoomController
	oidcController *controllers.OidcController
	userController *controllers.UserController
	configuration  config.IConfiguration
	jwtProvider    contracts.IJwtProvider
}

func NewServerDependencies(db *sqlx.DB, configuration config.IConfiguration) *ServerDependencies {
	encrypter := authentication.NewEncrypter(configuration)
	jwtProvider := authentication.NewJwtProvider(configuration)
	googleOidcService := authentication.NewGoogleOidcService(configuration)

	unitOfWork := persistance.NewUnitOfWork(db)
	userRepository := persistance.NewUserRepository()
	roomRepository := persistance.NewRoomRepository()
	externalCredentialsRepository := persistance.NewExternalCredentialsRepository(encrypter)
	refreshTokenRepository := persistance.NewRefreshTokenRepository(encrypter)

	registerUserCommandHandler := register.NewRegisterUserCommandHandler(userRepository,
		unitOfWork,
		externalCredentialsRepository,
		jwtProvider,
		refreshTokenRepository,
		encrypter,
	)
	loginUserCommandHandler := login.NewLoginUserCommandHandler(unitOfWork,
		userRepository,
		encrypter,
		jwtProvider,
		refreshTokenRepository,
		externalCredentialsRepository)
	getUserDataQueryHandler := data.NewUserQueryHandler(unitOfWork,
		userRepository)

	oidcAuthenticationService := services.NewOidcAuthenticationService(googleOidcService,
		userRepository,
		unitOfWork,
		registerUserCommandHandler,
		loginUserCommandHandler)

	oidcController := controllers.NewOidcController(configuration,
		oidcAuthenticationService,
		googleOidcService)
	joinRoomCommandHandler := join.NewRoomCommandHandler(roomRepository, unitOfWork)
	roomController := controllers.NewRoomController(joinRoomCommandHandler)
	userController := controllers.NewUserController(getUserDataQueryHandler)

	return &ServerDependencies{
		roomController: roomController,
		oidcController: oidcController,
		userController: userController,
		configuration:  configuration,
		jwtProvider:    jwtProvider,
	}
}
func (sd ServerDependencies) RoomController() *controllers.RoomController {
	return sd.roomController
}
func (sd ServerDependencies) OidcController() *controllers.OidcController {
	return sd.oidcController
}
func (sd ServerDependencies) UserController() *controllers.UserController {
	return sd.userController
}
func (sd ServerDependencies) Configuration() config.IConfiguration {
	return sd.configuration
}
func (sd ServerDependencies) JwtProvider() contracts.IJwtProvider {
	return sd.jwtProvider
}
