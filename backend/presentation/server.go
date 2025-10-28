package presentation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"xsedox.com/main/config"
	"xsedox.com/main/initialization"
	customMiddleware "xsedox.com/main/presentation/custom_middleware"
)

type Server struct {
	router *chi.Mux
}

func NewServer(dependencies *initialization.ServerDependencies) Server {
	router := chi.NewRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{dependencies.Configuration().Authentication().ClientOrigin},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: dependencies.Configuration().IsDevelopment(),
	})
	jwtAuthMiddleware := customMiddleware.NewJwtAuthentication(dependencies.Configuration(), dependencies.JwtProvider())

	router.Use(corsMiddleware.Handler,
		middleware.Logger,
		middleware.Recoverer)

	swaggerDocUrl := fmt.Sprintf("%v:%v/api/swagger/doc.json",
		dependencies.Configuration().Server().Host,
		dependencies.Configuration().Server().Port)
	router.Get("/api/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerDocUrl)))

	apiV1 := chi.NewRouter()
	securedRouter := chi.NewRouter()
	securedRouter.Use(jwtAuthMiddleware.Next)

	roomRouter := chi.NewRouter()
	roomRouter.Post("/", dependencies.RoomController().CreateRoom)

	userRouter := chi.NewRouter()

	oidcRouter := chi.NewRouter()
	oidcRouter.Get("/google/signin-oidc", dependencies.OidcController().HandleLoginWithGoogle)
	oidcRouter.Get("/google/callback", dependencies.OidcController().HandleGoogleCallback)

	securedRouter.Mount("/room", roomRouter)
	securedRouter.Mount("/login", userRouter)

	apiV1.Mount("/auth", oidcRouter)
	apiV1.Mount("/", securedRouter)
	router.Mount("/api/v1", apiV1)

	return Server{
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
