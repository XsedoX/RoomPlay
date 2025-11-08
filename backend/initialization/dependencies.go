package initialization

import (
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/create_command"
	"xsedox.com/main/application/services"
	"xsedox.com/main/application/user/data_query"
	"xsedox.com/main/application/user/login_command"
	"xsedox.com/main/application/user/login_refresh_token"
	"xsedox.com/main/application/user/logout_command"
	"xsedox.com/main/application/user/register_command"
	"xsedox.com/main/config"
	"xsedox.com/main/infrastructure/authentication"
	"xsedox.com/main/infrastructure/persistance"
	"xsedox.com/main/presentation/controllers"
)

type ServerDependencies struct {
	roomController           *controllers.RoomController
	oidcController           *controllers.OidcController
	userController           *controllers.UserController
	authenticationController *controllers.AuthenticationController
	configuration            config.IConfiguration
	jwtProvider              contracts.IJwtProvider
}

func NewServerDependencies(db *sqlx.DB, configuration config.IConfiguration) *ServerDependencies {
	encrypter := authentication.NewEncrypter(configuration)
	jwtProvider := authentication.NewJwtProvider(configuration)
	googleOidcService := authentication.NewGoogleOidcService(configuration)
	unitOfWork := persistance.NewUnitOfWork(db)

	userRepository := persistance.NewUserRepository()
	externalCredentialsRepository := persistance.NewExternalCredentialsRepository(encrypter)
	refreshTokenRepository := persistance.NewRefreshTokenRepository(encrypter)
	roomRepository := persistance.NewRoomRepository(encrypter)

	registerUserCommandHandler := register.NewRegisterUserCommandHandler(userRepository,
		unitOfWork,
		externalCredentialsRepository,
		jwtProvider,
		refreshTokenRepository,
		encrypter,
	)
	loginUserCommandHandler := login_command.NewLoginUserCommandHandler(unitOfWork,
		userRepository,
		encrypter,
		jwtProvider,
		refreshTokenRepository,
		externalCredentialsRepository)
	getUserDataQueryHandler := data_query.NewUserDataQueryHandler(unitOfWork,
		userRepository)

	loginRefreshTokenCommandHandler := login_refresh_token.NewLoginRefreshTokenCommandHandler(refreshTokenRepository,
		unitOfWork,
		encrypter,
		jwtProvider,
		userRepository)
	logoutRefreshTokenCommandHandler := logout_command.NewLogoutRefreshTokenCommandHandler(refreshTokenRepository, unitOfWork)

	createRoomCommandHandler := create_command.NewCreateRoomCommandHandler(roomRepository,
		unitOfWork,
		encrypter,
	)

	oidcAuthenticationService := services.NewOidcAuthenticationService(googleOidcService,
		userRepository,
		unitOfWork,
		registerUserCommandHandler,
		loginUserCommandHandler)

	oidcController := controllers.NewOidcController(configuration,
		oidcAuthenticationService,
		googleOidcService)
	userController := controllers.NewUserController(getUserDataQueryHandler)
	authenticationController := controllers.NewAuthenticationController(loginRefreshTokenCommandHandler,
		configuration,
		logoutRefreshTokenCommandHandler)
	roomController := controllers.NewRoomController(createRoomCommandHandler)

	return &ServerDependencies{
		oidcController:           oidcController,
		userController:           userController,
		authenticationController: authenticationController,
		roomController:           roomController,
		configuration:            configuration,
		jwtProvider:              jwtProvider,
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
func (sd ServerDependencies) AuthenticationController() *controllers.AuthenticationController {
	return sd.authenticationController
}
func (sd ServerDependencies) Configuration() config.IConfiguration {
	return sd.configuration
}
func (sd ServerDependencies) JwtProvider() contracts.IJwtProvider {
	return sd.jwtProvider
}
