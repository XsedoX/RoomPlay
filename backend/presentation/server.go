package presentation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"xsedox.com/main/config"
	"xsedox.com/main/presentation/handlers"
	customMiddleware "xsedox.com/main/presentation/middleware"
)

type ServerDependencies struct {
	RoomHandler   *handlers.RoomHandler
	UserHandler   *handlers.UserHandler
	Configuration *config.Configuration
}

type Server struct {
	roomHandler *handlers.RoomHandler
	userHandler *handlers.UserHandler
	router      *chi.Mux
	config      *config.Configuration
}

func NewServer(dependencies *ServerDependencies) Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	swaggerDocUrl := fmt.Sprintf("%v:%v/api/swagger/doc.json", dependencies.Configuration.Server.Host, dependencies.Configuration.Server.Port)
	router.Get("/api/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(swaggerDocUrl)))

	router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(rg chi.Router) {
			if dependencies.Configuration.IsProduction() {
				rg.Use(func(next http.Handler) http.Handler {
					return customMiddleware.Authentication(next, dependencies.Configuration)
				})
			}
			rg.Route("/room", func(routeRoom chi.Router) {
				routeRoom.Post("/", dependencies.RoomHandler.CreateRoom)
			})
		})
		r.Route("/user", func(routeUser chi.Router) {
			routeUser.Post("/", dependencies.UserHandler.LoginUser)
		})
	})

	return Server{
		roomHandler: dependencies.RoomHandler,
		router:      router,
		config:      dependencies.Configuration,
	}
}
func (s *Server) Start() {
	log.Printf("Starting API server on :%v", s.config.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", s.config.Server.Port), s.router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
