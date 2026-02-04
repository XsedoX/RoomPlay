package initialize_dependencies

import (
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_jwt_provider"
	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command_handler"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_handler"
	"github.com/XsedoX/RoomPlay/application/room/get_user_room_membership/get_user_room_membership_query_handler"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command_handler"
	"github.com/XsedoX/RoomPlay/application/room/leave_room/leave_room_command_handler"
	"github.com/XsedoX/RoomPlay/application/services/oidc_authentication_service"
	"github.com/XsedoX/RoomPlay/application/user/get_user/get_user_query_handler"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command_handler"
	"github.com/XsedoX/RoomPlay/application/user/login_user_refresh_token/login_user_refresh_token_command_handler"
	"github.com/XsedoX/RoomPlay/application/user/logout_user/logout_user_command_handler"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command_handler"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/google_oidc_service"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/jwt_provider"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/external_credentials_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/room/room_repository"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/unit_of_work"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/user/user_repository"
	"github.com/XsedoX/RoomPlay/presentation/controllers/authentication_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/google_oidc_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/user_controller"
	"github.com/jmoiron/sqlx"
)

type ServerDependencies struct {
	roomController           *room_controller.RoomController
	oidcController           *google_oidc_controller.GoogleOidcController
	userController           *user_controller.UserController
	authenticationController *authentication_controller.AuthenticationController
	configuration            config.IConfiguration
	jwtProvider              i_jwt_provider.IJwtProvider
}

func NewServerDependencies(db *sqlx.DB, configuration config.IConfiguration) *ServerDependencies {
	encrypter := encryper.NewEncrypter(configuration)
	jwtProvider := jwt_provider.NewJwtProvider(configuration)
	googleOidcService := google_oidc_service.NewGoogleOidcService(configuration)
	unitOfWork := unit_of_work.NewUnitOfWork(db)

	userRepository := user_repository.NewUserRepository()
	externalCredentialsRepository := external_credentials_repository.NewExternalCredentialsRepository(encrypter)
	refreshTokenRepository := internal_credentials_repository.NewInternalCredentialsRepository(encrypter)
	roomRepository := room_repository.NewRoomRepository(encrypter)

	registerUserCommandHandler := register_user_command_handler.NewRegisterUserCommandHandler(userRepository,
		unitOfWork,
		externalCredentialsRepository,
		jwtProvider,
		refreshTokenRepository,
		encrypter,
	)
	loginUserCommandHandler := login_user_command_handler.NewLoginUserCommandHandler(unitOfWork,
		userRepository,
		encrypter,
		jwtProvider,
		refreshTokenRepository,
		externalCredentialsRepository)
	getUserDataQueryHandler := get_user_query_handler.NewGetUserQueryHandler(unitOfWork,
		userRepository)
	loginRefreshTokenCommandHandler := login_user_refresh_token_command_handler.NewLoginUserRefreshTokenCommandHandler(refreshTokenRepository,
		unitOfWork,
		encrypter,
		jwtProvider,
		userRepository)
	logoutRefreshTokenCommandHandler := logout_user_command_handler.NewLogoutUserCommandHandler(refreshTokenRepository, unitOfWork)
	getUserRoomMembershipQueryHandler := get_user_room_membership_query_handler.NewGetUserRoomMembershipQueryHandler(roomRepository,
		unitOfWork)

	createRoomCommandHandler := create_room_command_handler.NewCreateRoomCommandHandler(roomRepository,
		unitOfWork,
		encrypter,
	)
	leaveRoomCommandHandler := leave_room_command_handler.NewLeaveRoomCommandHandler(roomRepository,
		unitOfWork,
	)
	getRoomQueryHandler := get_room_query_handler.NewGetRoomQueryHandler(unitOfWork,
		roomRepository,
	)
	joinRoomPasswordCommandHandler := join_room_password_command_handler.NewJoinRoomPasswordCommandHandler(roomRepository, unitOfWork)

	oidcAuthenticationService := oidc_authentication_service.NewOidcAuthenticationService(googleOidcService,
		userRepository,
		unitOfWork,
		registerUserCommandHandler,
		loginUserCommandHandler)

	oidcController := google_oidc_controller.NewOidcController(configuration,
		oidcAuthenticationService,
		googleOidcService)
	userController := user_controller.NewUserController(getUserDataQueryHandler)
	authenticationController := authentication_controller.NewAuthenticationController(loginRefreshTokenCommandHandler,
		configuration,
		logoutRefreshTokenCommandHandler)
	roomController := room_controller.NewRoomController(createRoomCommandHandler,
		getRoomQueryHandler,
		getUserRoomMembershipQueryHandler,
		leaveRoomCommandHandler,
		joinRoomPasswordCommandHandler)

	return &ServerDependencies{
		oidcController:           oidcController,
		userController:           userController,
		authenticationController: authenticationController,
		roomController:           roomController,
		configuration:            configuration,
		jwtProvider:              jwtProvider,
	}
}

func (sd ServerDependencies) RoomController() *room_controller.RoomController {
	return sd.roomController
}

func (sd ServerDependencies) OidcController() *google_oidc_controller.GoogleOidcController {
	return sd.oidcController
}

func (sd ServerDependencies) UserController() *user_controller.UserController {
	return sd.userController
}

func (sd ServerDependencies) AuthenticationController() *authentication_controller.AuthenticationController {
	return sd.authenticationController
}

func (sd ServerDependencies) Configuration() config.IConfiguration {
	return sd.configuration
}

func (sd ServerDependencies) JwtProvider() i_jwt_provider.IJwtProvider {
	return sd.jwtProvider
}
