package presentation_dependencies

import (
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/presentation/application_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/controllers/authentication_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/google_oidc_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/user_controller"
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
)

type PresentationDependencies struct {
	roomController           *room_controller.RoomController
	oidcController           *google_oidc_controller.GoogleOidcController
	userController           *user_controller.UserController
	authenticationController *authentication_controller.AuthenticationController
}

func ConstructPresentationDependencies(
	configuration config.IConfiguration,
	applicationDependencies *application_dependencies.ApplicationDependencies,
	infrastructureDependencies *infrastructure_dependencies.InfrastructureDependencies,
) *PresentationDependencies {
	oidcAuthenticationService := applicationDependencies.OidcAuthenticationService
	googleOidcService := infrastructureDependencies.GoogleOidcService
	getUserDataQueryHandler := applicationDependencies.GetUserDataQueryHandler
	loginRefreshTokenCommandHandler := applicationDependencies.LoginRefreshTokenCommandHandler
	logoutRefreshTokenCommandHandler := applicationDependencies.LogoutRefreshTokenCommandHandler

	createRoomCommandHandler := applicationDependencies.CreateRoomCommandHandler
	getRoomQueryHandler := applicationDependencies.GetRoomQueryHandler
	getUserRoomMembershipQueryHandler := applicationDependencies.GetUserRoomMembershipQueryHandler
	leaveRoomCommandHandler := applicationDependencies.LeaveRoomCommandHandler
	joinRoomPasswordCommandHandler := applicationDependencies.JoinRoomPasswordCommandHandler

	oidcController := google_oidc_controller.NewOidcController(
		configuration,
		oidcAuthenticationService,
		googleOidcService,
	)
	userController := user_controller.NewUserController(getUserDataQueryHandler)

	authenticationController := authentication_controller.NewAuthenticationController(
		loginRefreshTokenCommandHandler,
		configuration,
		logoutRefreshTokenCommandHandler)

	roomController := room_controller.NewRoomController(
		createRoomCommandHandler,
		getRoomQueryHandler,
		getUserRoomMembershipQueryHandler,
		leaveRoomCommandHandler,
		joinRoomPasswordCommandHandler)

	return &PresentationDependencies{
		oidcController:           oidcController,
		userController:           userController,
		authenticationController: authenticationController,
		roomController:           roomController,
	}
}

func (pd *PresentationDependencies) RoomController() *room_controller.RoomController {
	return pd.roomController
}

func (pd *PresentationDependencies) OidcController() *google_oidc_controller.GoogleOidcController {
	return pd.oidcController
}

func (pd *PresentationDependencies) UserController() *user_controller.UserController {
	return pd.userController
}

func (pd *PresentationDependencies) AuthenticationController() *authentication_controller.AuthenticationController {
	return pd.authenticationController
}
