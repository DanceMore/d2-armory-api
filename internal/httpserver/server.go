package httpserver

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// Server is the HTTP server listener.
type Server struct {
	encoder           *encoder
	listener          net.Listener
	addr              string
	characterService  characterService
	statisticsService statisticsService
	credentials       map[string]string
	corsEnabled       bool
	loggingEnabled    bool
}

// Open will open a tcp listener to serve http requests.
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = ln

	// Create an http server.
	server := http.Server{
		Handler:     http.TimeoutHandler(s.Handler(), (2 * time.Second), "connection timeout"),
		ReadTimeout: 5 * time.Second,
	}

	log.Println("starting HTTP server on:", s.addr)

	return server.Serve(s.listener)
}

// Handler will setup a router that implements the http.Handler interface.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	if s.loggingEnabled {
		// Middleware for logging requests.
		r.Use(middleware.Logger)
	}

	if s.corsEnabled {
		cors := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})

		r.Use(cors.Handler)
	}

	r.Route("/health", newHealthHandler().Routes)
	r.Route("/api/v1/characters", newCharacterHandler(s.encoder, s.characterService).Routes)
	r.Route("/api/v1/statistics", newStatisticsHandler(s.encoder, s.statisticsService, s.credentials).Routes)

	// Deprecated handler, supported for consumers who rely on it.
	r.Route("/retrieving/v1/character", newCharacterHandler(s.encoder, s.characterService).Routes)

        // new functionality not used at SlashDiablo
        // TODO: submit it upstream ;)
        // TODO: make it conditional for upstream ...?
        r.Route("/list-characters", newListCharactersHandler().Routes)

	return r
}

// NewServer returns a new server with all dependencies.
func NewServer(addr string, characterService characterService, statisticsService statisticsService, credentials map[string]string, corsEnabled bool, loggingEnabled bool) *Server {
	return &Server{
		addr:              addr,
		encoder:           newEncoder(),
		characterService:  characterService,
		statisticsService: statisticsService,
		credentials:       credentials,
		corsEnabled:       corsEnabled,
		loggingEnabled:    loggingEnabled,
	}
}
