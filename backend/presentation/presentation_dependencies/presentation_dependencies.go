package presentation_dependencies

import (
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/domain/room/events"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_handlers/song_added_client_message_handler"
	"github.com/XsedoX/RoomPlay/infrastructure/event_handlers/song_enqueued_websocket_event"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/main_hub"
	"github.com/XsedoX/RoomPlay/presentation/application_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/controllers/authentication_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/google_oidc_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/song_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/user_controller"
	"github.com/XsedoX/RoomPlay/presentation/infrastructure_dependencies"
)

type PresentationDependencies struct {
	roomController           *room_controller.RoomController
	oidcController           *google_oidc_controller.GoogleOidcController
	userController           *user_controller.UserController
	authenticationController *authentication_controller.AuthenticationController
	songController           *song_controller.SongController
}

func ConstructPresentationDependencies(
	configuration config.IConfiguration,
	applicationDependencies *application_dependencies.ApplicationDependencies,
	infrastructureDependencies *infrastructure_dependencies.InfrastructureDependencies,
) *PresentationDependencies {
	googleOidcService := infrastructureDependencies.GoogleOidcService
	oidcAuthenticationService := applicationDependencies.OidcAuthenticationService
	oidcController := google_oidc_controller.NewOidcController(
		configuration,
		oidcAuthenticationService,
		googleOidcService,
	)

	mainHub := main_hub.NewHub(infrastructureDependencies.ApplicationContext)

	songEnqueuedWebsocketEventHandler := song_enqueued_websocket_event.NewSongEnqueuedWebsocketEventHandler(
		mainHub,
		infrastructureDependencies.RoomRepository,
		infrastructureDependencies.UnitOfWork,
		infrastructureDependencies.ApplicationContext,
	)
	infrastructureDependencies.DomainEventPublisher.Register(
		events.SongEnqueuedEventName,
		songEnqueuedWebsocketEventHandler,
	)
	enqueueSongCommandHandler := applicationDependencies.EnqueueSongCommandHandler
	clientMessageHandler := song_added_client_message_handler.NewSongAddedClientMessageHandler(
		enqueueSongCommandHandler,
	)
	infrastructureDependencies.ClientMessagePublisher.RegisterHandler(
		song_added_client_message_handler.SongAddedClientMessageName,
		clientMessageHandler,
	)

	getUserDataQueryHandler := applicationDependencies.GetUserDataQueryHandler
	userController := user_controller.NewUserController(getUserDataQueryHandler)

	loginRefreshTokenCommandHandler := applicationDependencies.LoginRefreshTokenCommandHandler
	logoutRefreshTokenCommandHandler := applicationDependencies.LogoutRefreshTokenCommandHandler
	authenticationController := authentication_controller.NewAuthenticationController(
		loginRefreshTokenCommandHandler,
		configuration,
		logoutRefreshTokenCommandHandler)

	joinRoomPasswordCommandHandler := applicationDependencies.JoinRoomPasswordCommandHandler
	leaveRoomCommandHandler := applicationDependencies.LeaveRoomCommandHandler
	getRoomQueryHandler := applicationDependencies.GetRoomQueryHandler
	createRoomCommandHandler := applicationDependencies.CreateRoomCommandHandler
	getUserRoomMembershipQueryHandler := applicationDependencies.GetUserRoomMembershipQueryHandler
	roomController := room_controller.NewRoomController(
		createRoomCommandHandler,
		getRoomQueryHandler,
		getUserRoomMembershipQueryHandler,
		leaveRoomCommandHandler,
		joinRoomPasswordCommandHandler,
		mainHub,
		infrastructureDependencies.ClientMessagePublisher,
	)

	searchSongQueryHandler := applicationDependencies.SearchSongQueryHandler
	songController := song_controller.NewSongController(
		searchSongQueryHandler,
	)

	return &PresentationDependencies{
		oidcController:           oidcController,
		userController:           userController,
		authenticationController: authenticationController,
		songController:           songController,
		roomController:           roomController,
	}
}

func (pd *PresentationDependencies) SongController() *song_controller.SongController {
	return pd.songController
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
