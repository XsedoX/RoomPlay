package presentation

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	roomHandler *RoomHandler
	router      *chi.Mux
}

func NewServer(roomHandler *RoomHandler) Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/room", roomHandler.CreateRoom)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7865/swagger/doc.json")))

	return Server{
		roomHandler: roomHandler,
		router:      router,
	}
}
func (s *Server) Start() {
	log.Println("Starting API server on :7865")
	if err := http.ListenAndServe(":7865", s.router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
