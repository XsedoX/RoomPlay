package presentation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"xsedox.com/main/config"
	"xsedox.com/main/initialization"
	customMiddleware "xsedox.com/main/presentation/custom_middleware"
)

type Server struct {
	router *chi.Mux
}

func NewServer(dependencies *initialization.ServerDependencies) *Server {
	router := chi.NewRouter()
	customCors := customMiddleware.NewCustomCors(dependencies.Configuration())
	jwtAuthMiddleware := customMiddleware.NewCookieJwtAuthentication(dependencies.Configuration(), dependencies.JwtProvider())

	router.Use(customCors.CorsHandler(),
		middleware.Logger,
		middleware.Recoverer)

	swaggerDocUrl := fmt.Sprintf("%v:%v/api/swagger/doc.json",
		(dependencies.Configuration()).Server().Host,
		(dependencies.Configuration()).Server().Port)
	router.Get("/api/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerDocUrl)))

	apiV1 := chi.NewRouter()

	// Public routes
	apiV1.Get("/auth/google/signin-oidc", dependencies.OidcController().HandleLoginWithGoogle)
	apiV1.Get("/auth/google/callback", dependencies.OidcController().HandleGoogleCallback)
	apiV1.Post("/auth/refresh-token", dependencies.AuthenticationController().RefreshToken)

	// Secured routes
	apiV1.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware.Next)

		r.Post("/auth/logout", dependencies.AuthenticationController().Logout)

		r.Route("/room", func(r chi.Router) {
			r.Post("/", dependencies.RoomController().CreateRoom)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", dependencies.UserController().GetUserData)
		})
	})

	router.Mount("/api/v1", apiV1)

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
