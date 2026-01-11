package initialization

import (
	"github.com/jmoiron/sqlx"
	"github.com/XsedoX/RoomPlay/application/contracts"
	"github.com/XsedoX/RoomPlay/application/room/create_room"
	"github.com/XsedoX/RoomPlay/application/room/get_room"
	"github.com/XsedoX/RoomPlay/application/room/get_user_room_membership"
	"github.com/XsedoX/RoomPlay/application/room/leave_room"
	"github.com/XsedoX/RoomPlay/application/services"
	"github.com/XsedoX/RoomPlay/application/user/get_user"
	"github.com/XsedoX/RoomPlay/application/user/login_user"
	"github.com/XsedoX/RoomPlay/application/user/login_user_refresh_token"
	"github.com/XsedoX/RoomPlay/application/user/logout_user"
	"github.com/XsedoX/RoomPlay/application/user/register_user"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance"
	"github.com/XsedoX/RoomPlay/presentation/controllers"
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

	registerUserCommandHandler := register_user.NewRegisterUserCommandHandler(userRepository,
		unitOfWork,
		externalCredentialsRepository,
		jwtProvider,
		refreshTokenRepository,
		encrypter,
	)
	loginUserCommandHandler := login_user.NewLoginUserCommandHandler(unitOfWork,
		userRepository,
		encrypter,
		jwtProvider,
		refreshTokenRepository,
		externalCredentialsRepository)
	getUserDataQueryHandler := get_user.NewGetUserQueryHandler(unitOfWork,
		userRepository)

	loginRefreshTokenCommandHandler := login_user_refresh_token.NewLoginUserRefreshTokenCommandHandler(refreshTokenRepository,
		unitOfWork,
		encrypter,
		jwtProvider,
		userRepository)
	logoutRefreshTokenCommandHandler := logout_user.NewLogoutUserCommandHandler(refreshTokenRepository, unitOfWork)

	createRoomCommandHandler := create_room.NewCreateRoomCommandHandler(roomRepository,
		unitOfWork,
		encrypter,
	)
	leaveRoomCommandHandler := leave_room.NewLeaveRoomCommandHandler(roomRepository,
		unitOfWork,
	)
	getRoomQueryHandler := get_room.NewGetRoomQueryHandler(unitOfWork,
		roomRepository,
	)
	getUserRoomMembershipQueryHandler := get_user_room_membership.NewGetUserRoomMembershipQueryHandler(roomRepository,
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
