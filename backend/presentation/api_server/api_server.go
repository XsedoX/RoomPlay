package api_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/song_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/user_controller"
	"github.com/XsedoX/RoomPlay/presentation/custom_middleware/cookie_jwt_authentication_middleware"
	"github.com/XsedoX/RoomPlay/presentation/custom_middleware/cors_middleware"
	"github.com/XsedoX/RoomPlay/presentation/custom_middleware/security_headers_middleware"
	"github.com/XsedoX/RoomPlay/presentation/initialize_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router     *chi.Mux
	appContext context.Context
}

func NewServer(dependencies *initialize_dependencies.ServerDependencies, configuration config.IConfiguration, customMiddlewares ...func(http.Handler) http.Handler) *Server {
	router := chi.NewRouter()
	customCors := cors_middleware.NewCustomCors(configuration)
	jwtAuthMiddleware := cookie_jwt_authentication_middleware.NewCookieJwtAuthentication(configuration, dependencies.InfrastructureDependencies.JwtProvider)
	securityHeadersMiddleware := security_headers_middleware.NewSecurityHeaders(configuration)
	// Apply custom middlewares
	router.Use(customMiddlewares...)

	if !configuration.IsTesting() {
		router.Use(securityHeadersMiddleware.Next,
			customCors.CorsHandler(),
			middleware.Logger,
			middleware.Recoverer)
	}

	apiV1 := chi.NewRouter()

	// Public routes
	apiV1.Get("/auth/google/signin-oidc", dependencies.PresentationDependencies.OidcController().HandleLoginWithGoogle)
	apiV1.Get("/auth/google/callback", dependencies.PresentationDependencies.OidcController().HandleGoogleCallback)
	apiV1.Post(constants.RefreshTokenPath, dependencies.PresentationDependencies.AuthenticationController().RefreshToken)

	// Secured routes
	apiV1.Group(func(r chi.Router) {
		if !configuration.IsTesting() {
			r.Use(jwtAuthMiddleware.Next)
		}

		r.Post(constants.LogoutPath, dependencies.PresentationDependencies.AuthenticationController().Logout)

		r.Route(room_controller.RoomBasePath, func(r chi.Router) {
			r.Post("/", dependencies.PresentationDependencies.RoomController().CreateRoom)
			r.Get("/", dependencies.PresentationDependencies.RoomController().GetRoom)
			r.Delete("/", dependencies.PresentationDependencies.RoomController().LeaveRoom)
			r.Get(room_controller.RoomMembershipBasePath, dependencies.PresentationDependencies.RoomController().CheckUserRoomMembership)
			r.Put(room_controller.JoinRoomPasswordPath, dependencies.PresentationDependencies.RoomController().JoinRoomPassword)
		})

		r.Route(song_controller.SongBasePath, func(r chi.Router) {
			r.Get(song_controller.SearchSongPath, dependencies.PresentationDependencies.SongController().SearchSongsByQuery)
		})

		r.Route(user_controller.UserBasePath, func(r chi.Router) {
			r.Get("/", dependencies.PresentationDependencies.UserController().GetUserData)
		})
	})

	router.Mount(constants.ApiBasePath, apiV1)

	return &Server{
		router:     router,
		appContext: dependencies.InfrastructureDependencies.ApplicationContext,
	}
}

func (s *Server) Start(configuration config.IConfiguration) {
	log.Printf("Starting API server on :%v", configuration.Server().Port)
	// err := http.ListenAndServeTLS(fmt.Sprintf(":%v", configuration.Server().Port), "./certificates/server.crt", "./certificates/server.key", s.router)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", configuration.Server().Port),
		Handler: s.router,
	}
	go func() {
		<-s.appContext.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) UpdateRouter(newRouter *chi.Mux) {
	s.router = newRouter
}
