package webserver

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.Handler
	WebServerPort string
}

func NewWebServer(servePort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.Handler),
		WebServerPort: servePort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.Handler) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	// Middlewares globais
	s.Router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
	)

	for path, handler := range s.Handlers {
		s.Router.Mount(path, handler)
		log.Printf("Route mounted: %s", path)
	}

	log.Printf("Server running on port %s", s.WebServerPort)
	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}
