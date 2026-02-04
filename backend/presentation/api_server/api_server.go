package api_server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/initialization/initialize_dependencies"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/controllers/user_controller"
	"github.com/XsedoX/RoomPlay/presentation/custom_middleware/cookie_jwt_authentication_middleware"
	"github.com/XsedoX/RoomPlay/presentation/custom_middleware/cors_middleware"
	"github.com/XsedoX/RoomPlay/presentation/custom_middleware/security_headers_middleware"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	router *chi.Mux
}

func NewServer(dependencies *initialize_dependencies.ServerDependencies, customMiddlewares ...func(http.Handler) http.Handler) *Server {
	router := chi.NewRouter()
	customCors := cors_middleware.NewCustomCors(dependencies.Configuration())
	jwtAuthMiddleware := cookie_jwt_authentication_middleware.NewCookieJwtAuthentication(dependencies.Configuration(), dependencies.JwtProvider())
	securityHeadersMiddleware := security_headers_middleware.NewSecurityHeaders(dependencies.Configuration())
	// Apply custom middlewares
	router.Use(customMiddlewares...)

	if !dependencies.Configuration().IsTesting() {
		router.Use(securityHeadersMiddleware.Next,
			customCors.CorsHandler(),
			middleware.Logger,
			middleware.Recoverer)

		swaggerDocUrl := fmt.Sprintf("%v:%v/api/swagger/doc.json",
			(dependencies.Configuration()).Server().Host,
			(dependencies.Configuration()).Server().Port)
		router.Get("/api/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerDocUrl)))
	}

	apiV1 := chi.NewRouter()

	// Public routes
	apiV1.Get("/auth/google/signin-oidc", dependencies.OidcController().HandleLoginWithGoogle)
	apiV1.Get("/auth/google/callback", dependencies.OidcController().HandleGoogleCallback)
	apiV1.Post(constants.AuthBasePath+constants.RefreshTokenPath, dependencies.AuthenticationController().RefreshToken)

	// Secured routes
	apiV1.Group(func(r chi.Router) {
		if !dependencies.Configuration().IsTesting() {
			r.Use(jwtAuthMiddleware.Next)
		}

		r.Post(constants.AuthBasePath+constants.LogoutPath, dependencies.AuthenticationController().Logout)

		r.Route(room_controller.RoomBasePath, func(r chi.Router) {
			r.Post("/", dependencies.RoomController().CreateRoom)
			r.Get("/", dependencies.RoomController().GetRoom)
			r.Delete("/", dependencies.RoomController().LeaveRoom)
			r.Get(room_controller.RoomMembershipBasePath, dependencies.RoomController().CheckUserRoomMembership)
			r.Put(room_controller.JoinRoomPasswordPath, dependencies.RoomController().JoinRoomPassword)
		})

		r.Route(user_controller.UserBasePath, func(r chi.Router) {
			r.Get("/", dependencies.UserController().GetUserData)
		})
	})

	router.Mount(constants.ApiBasePath, apiV1)

	return &Server{
		router: router,
	}
}

func (s *Server) Start(configuration config.IConfiguration) {
	log.Printf("Starting API server on :%v", configuration.Server().Port)
	// err := http.ListenAndServeTLS(fmt.Sprintf(":%v", configuration.Server().Port), "./certificates/server.crt", "./certificates/server.key", s.router)
	err := http.ListenAndServe(fmt.Sprintf(":%v", configuration.Server().Port), s.router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) UpdateRouter(newRouter *chi.Mux) {
	s.router = newRouter
}
