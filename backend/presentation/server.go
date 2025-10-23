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
	"xsedox.com/main/presentation/handlers"
)

type Server struct {
	roomHandler *handlers.RoomHandler
	userHandler *handlers.UserHandler
	oidcHandler *handlers.OidcHandler
	router      *chi.Mux
}

func NewServer(dependencies *initialization.ServerDependencies) Server {
	router := chi.NewRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{config.Config.Authentication.ClientOrigin},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: config.Config.IsDevelopment(),
	})

	router.Use(corsMiddleware.Handler,
		middleware.Logger,
		middleware.Recoverer)

	swaggerDocUrl := fmt.Sprintf("%v:%v/api/swagger/doc.json", config.Config.Server.Host, config.Config.Server.Port)
	router.Get("/api/swagger/*", httpSwagger.Handler(httpSwagger.URL(swaggerDocUrl)))

	apiV1 := chi.NewRouter()
	securedRouter := chi.NewRouter()
	securedRouter.Use(customMiddleware.JwtAuthentication)

	roomRouter := chi.NewRouter()
	roomRouter.Post("/", dependencies.RoomHandler().CreateRoom)

	userRouter := chi.NewRouter()
	userRouter.Post("/", dependencies.UserHandler().LoginUser)

	oidcRouter := chi.NewRouter()
	oidcRouter.Post("/google/signin-oidc", dependencies.OidcHandler().HandleLoginWithGoogle)

	securedRouter.Mount("/room", roomRouter)
	securedRouter.Mount("/user", userRouter)

	apiV1.Mount("/auth", oidcRouter)
	apiV1.Mount("/", securedRouter)
	router.Mount("/api/v1", apiV1)

	return Server{
		roomHandler: dependencies.RoomHandler(),
		userHandler: dependencies.UserHandler(),
		oidcHandler: dependencies.OidcHandler(),
		router:      router,
	}
}
func (s *Server) Start() {
	log.Printf("Starting API server on :%v", config.Config.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", config.Config.Server.Port), s.router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
