package presentation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/XsedoX/RoomPlay/config"
	"github.com/XsedoX/RoomPlay/initialization"
	"github.com/XsedoX/RoomPlay/presentation/controllers"
	customMiddleware "github.com/XsedoX/RoomPlay/presentation/custom_middleware"
	"github.com/XsedoX/RoomPlay/presentation/helpers"
)

type Server struct {
	router *chi.Mux
}

func NewServer(dependencies *initialization.ServerDependencies, customMiddlewares ...func(http.Handler) http.Handler) *Server {
	router := chi.NewRouter()
	customCors := customMiddleware.NewCustomCors(dependencies.Configuration())
	jwtAuthMiddleware := customMiddleware.NewCookieJwtAuthentication(dependencies.Configuration(), dependencies.JwtProvider())
	securityHeadersMiddleware := customMiddleware.NewSecurityHeaders(dependencies.Configuration())
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
	apiV1.Post("/auth/refresh-token", dependencies.AuthenticationController().RefreshToken)

	// Secured routes
	apiV1.Group(func(r chi.Router) {
		if !dependencies.Configuration().IsTesting() {
			r.Use(jwtAuthMiddleware.Next)
		}

		r.Post("/auth/logout", dependencies.AuthenticationController().Logout)

		r.Route(controllers.RoomBasePath, func(r chi.Router) {
			r.Post("/", dependencies.RoomController().CreateRoom)
			r.Get("/", dependencies.RoomController().GetRoom)
			r.Delete("/", dependencies.RoomController().LeaveRoom)
			r.Get(controllers.RoomMembershipBasePath, dependencies.RoomController().CheckUserRoomMembership)
		})

		r.Route(controllers.UserBasePath, func(r chi.Router) {
			r.Get("/", dependencies.UserController().GetUserData)
		})
	})

	router.Mount(helpers.ApiBasePath, apiV1)

	return &Server{
		router: router,
	}
}
func (s *Server) Start(configuration config.IConfiguration) {
	log.Printf("Starting API server on :%v", configuration.Server().Port)
	//err := http.ListenAndServeTLS(fmt.Sprintf(":%v", configuration.Server().Port), "./certificates/server.crt", "./certificates/server.key", s.router)
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
