package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	for path, handler := range s.Handlers {
		s.Router.Mount(path, handler)
		log.Printf("Route mounted: %s", path)
	}

	log.Printf("Server running on port %s", s.WebServerPort)
	http.ListenAndServe(":"+s.WebServerPort, s.Router)
}
