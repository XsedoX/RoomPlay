package application_dependencies

import (
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_command_handler"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_query_handler"
	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command"
	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command_handler"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_handler"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/application/room/get_user_room_membership/get_user_room_membership_query_handler"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command_handler"
	"github.com/XsedoX/RoomPlay/application/room/leave_room/leave_room_command"
	"github.com/XsedoX/RoomPlay/application/room/leave_room/leave_room_command_handler"
	"github.com/XsedoX/RoomPlay/application/services/oidc_authentication_service"
	"github.com/XsedoX/RoomPlay/application/services/services_contracts/i_oidc_authentication_service"
	get_user_data_query_handler "github.com/XsedoX/RoomPlay/application/user/get_user/get_user_query_handler"
	get_user_data_query_response "github.com/XsedoX/RoomPlay/application/user/get_user/get_user_query_response"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command_handler"
	"github.com/XsedoX/RoomPlay/application/user/login_user/login_user_command_response"
	"github.com/XsedoX/RoomPlay/application/user/login_user_refresh_token/login_user_refresh_token_command_handler"
	"github.com/XsedoX/RoomPlay/application/user/login_user_refresh_token/login_user_refresh_token_command_response"
	"github.com/XsedoX/RoomPlay/application/user/logout_user/logout_user_command"
	"github.com/XsedoX/RoomPlay/application/user/logout_user/logout_user_command_handler"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command_handler"
	"github.com/XsedoX/RoomPlay/application/user/register_user/register_user_command_response"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/persistance_dependencies"
)

type ApplicationDependencies struct {
	RegisterUserCommandHandler       i_command_handler.ICommandHandlerWithResponse[*register_user_command.RegisterUserCommand, *register_user_command_response.RegisterUserCommandResponse]
	LoginUserCommandHandler          i_command_handler.ICommandHandlerWithResponse[*login_user_command.LoginUserCommand, *login_user_command_response.LoginUserCommandResponse]
	LoginRefreshTokenCommandHandler  i_command_handler.ICommandHandlerWithResponse[*string, *login_user_refresh_token_command_response.LoginUserRefreshTokenCommandResponse]
	LogoutRefreshTokenCommandHandler i_command_handler.ICommandHandler[*logout_user_command.LogoutUserCommand]

	GetUserRoomMembershipQueryHandler i_query_handler.IQueryHandler[*bool]
	GetUserDataQueryHandler           i_query_handler.IQueryHandler[*get_user_data_query_response.GetUserDataQueryResponse]

	CreateRoomCommandHandler       i_command_handler.ICommandHandler[*create_room_command.CreateRoomCommand]
	LeaveRoomCommandHandler        i_command_handler.ICommandHandler[*leave_room_command.LeaveRoomCommand]
	JoinRoomPasswordCommandHandler i_command_handler.ICommandHandler[*join_room_password_command.JoinRoomPasswordCommand]

	GetRoomQueryHandler i_query_handler.IQueryHandler[*get_room_query_response.GetRoomQueryResponse]

	OidcAuthenticationService i_oidc_authentication_service.IOidcAuthenticationService
}

func ConstructApplicationDependencies(
	persistanceDependencies *persistance_dependencies.PersistanceDependencies,
	infrastructureDependencies *infrastructure_dependencies.InfrastructureDependencies,
	configuration config.IConfiguration,
) *ApplicationDependencies {
	googleOidcService := infrastructureDependencies.GoogleOidcService
	roomRepository := persistanceDependencies.RoomRepository
	userRepository := persistanceDependencies.UserRepository
	unitOfWork := persistanceDependencies.UnitOfWork
	externalCredentialsRepository := persistanceDependencies.ExternalCredentialsRepository
	internalCredentialsRepository := persistanceDependencies.InternalCredentialsRepository
	encrypter := infrastructureDependencies.Encrypter
	jwtProvider := infrastructureDependencies.JwtProvider

	registerUserCommandHandler := register_user_command_handler.NewRegisterUserCommandHandler(
		userRepository,
		unitOfWork,
		externalCredentialsRepository,
		jwtProvider,
		internalCredentialsRepository,
		encrypter,
	)
	loginUserCommandHandler := login_user_command_handler.NewLoginUserCommandHandler(
		unitOfWork,
		userRepository,
		encrypter,
		jwtProvider,
		internalCredentialsRepository,
		externalCredentialsRepository,
	)

	getUserDataQueryHandler := get_user_data_query_handler.NewGetUserQueryHandler(
		unitOfWork,
		userRepository,
	)

	loginRefreshTokenCommandHandler := login_user_refresh_token_command_handler.NewLoginUserRefreshTokenCommandHandler(
		internalCredentialsRepository,
		unitOfWork,
		encrypter,
		jwtProvider,
		userRepository,
	)

	logoutRefreshTokenCommandHandler := logout_user_command_handler.NewLogoutUserCommandHandler(
		internalCredentialsRepository,
		unitOfWork,
	)

	getUserRoomMembershipQueryHandler := get_user_room_membership_query_handler.NewGetUserRoomMembershipQueryHandler(
		roomRepository,
		unitOfWork,
	)

	createRoomCommandHandler := create_room_command_handler.NewCreateRoomCommandHandler(
		roomRepository,
		unitOfWork,
		encrypter,
	)

	leaveRoomCommandHandler := leave_room_command_handler.NewLeaveRoomCommandHandler(
		roomRepository,
		unitOfWork,
	)

	getRoomQueryHandler := get_room_query_handler.NewGetRoomQueryHandler(
		unitOfWork,
		roomRepository,
	)

	joinRoomPasswordCommandHandler := join_room_password_command_handler.NewJoinRoomPasswordCommandHandler(
		roomRepository,
		unitOfWork,
	)

	oidcAuthenticationService := oidc_authentication_service.NewOidcAuthenticationService(
		googleOidcService,
		userRepository,
		unitOfWork,
		registerUserCommandHandler,
		loginUserCommandHandler,
	)

	return &ApplicationDependencies{
		RegisterUserCommandHandler:        registerUserCommandHandler,
		LoginUserCommandHandler:           loginUserCommandHandler,
		GetUserDataQueryHandler:           getUserDataQueryHandler,
		LoginRefreshTokenCommandHandler:   loginRefreshTokenCommandHandler,
		LogoutRefreshTokenCommandHandler:  logoutRefreshTokenCommandHandler,
		GetUserRoomMembershipQueryHandler: getUserRoomMembershipQueryHandler,
		CreateRoomCommandHandler:          createRoomCommandHandler,
		LeaveRoomCommandHandler:           leaveRoomCommandHandler,
		GetRoomQueryHandler:               getRoomQueryHandler,
		JoinRoomPasswordCommandHandler:    joinRoomPasswordCommandHandler,
		OidcAuthenticationService:         oidcAuthenticationService,
	}
}
