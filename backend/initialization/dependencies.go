package initialization

import (
	"github.com/jmoiron/sqlx"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/room/create_room_command"
	"xsedox.com/main/application/room/get_room_query"
	"xsedox.com/main/application/room/get_user_room_membership_query"
	"xsedox.com/main/application/room/leave_room_command"
	"xsedox.com/main/application/services"
	"xsedox.com/main/application/user/get_user_query"
	"xsedox.com/main/application/user/login_user_command"
	"xsedox.com/main/application/user/login_user_refresh_token_command"
	"xsedox.com/main/application/user/logout_user_command"
	"xsedox.com/main/application/user/register_user_command"
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

	registerUserCommandHandler := register_user_command.NewRegisterUserCommandHandler(userRepository,
		unitOfWork,
		externalCredentialsRepository,
		jwtProvider,
		refreshTokenRepository,
		encrypter,
	)
	loginUserCommandHandler := login_user_command.NewLoginUserCommandHandler(unitOfWork,
		userRepository,
		encrypter,
		jwtProvider,
		refreshTokenRepository,
		externalCredentialsRepository)
	getUserDataQueryHandler := get_user_query.NewGetUserQueryHandler(unitOfWork,
		userRepository)

	loginRefreshTokenCommandHandler := login_user_refresh_token_command.NewLoginUserRefreshTokenCommandHandler(refreshTokenRepository,
		unitOfWork,
		encrypter,
		jwtProvider,
		userRepository)
	logoutRefreshTokenCommandHandler := logout_user_command.NewLogoutUserCommandHandler(refreshTokenRepository, unitOfWork)

	createRoomCommandHandler := create_room_command.NewCreateRoomCommandHandler(roomRepository,
		unitOfWork,
		encrypter,
	)
	leaveRoomCommandHandler := leave_room_command.NewLeaveRoomCommandHandler(userRepository,
		unitOfWork,
	)
	getRoomQueryHandler := get_room_query.NewGetRoomQueryHandler(unitOfWork,
		roomRepository,
	)
	getUserRoomMembershipQueryHandler := get_user_room_membership_query.NewGetUserRoomMembershipQueryHandler(roomRepository,
		unitOfWork)

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
	roomController := controllers.NewRoomController(createRoomCommandHandler,
		getRoomQueryHandler,
		getUserRoomMembershipQueryHandler,
		leaveRoomCommandHandler)

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
