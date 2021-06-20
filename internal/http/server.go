package http

import (
	"log"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/GregoryAlbouy/shrinker/pkg/httplog"
	"github.com/GregoryAlbouy/shrinker/pkg/queue"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type Server struct {
	*http.Server
	router *mux.Router
	Repository
	imageQueue queue.Producer
}

// Repository exposes the available operations to access the data layer.
// Operations are grouped under services.
type Repository struct {
	UserService internal.UserService
}

// NewServer returns a new instance of Server given configuration parameters.
func NewServer(addr string, repo Repository, q *amqp.Connection, verbose bool) (*Server, error) {
	prod, err := queue.NewProducer(q, verbose)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Server: &http.Server{Addr: addr},
		router: mux.NewRouter().StrictSlash(true),
		Repository: Repository{
			UserService: repo.UserService,
		},
		imageQueue: prod,
	}

	if verbose {
		s.router.Use(httplog.RequestLogger)
	}

	s.registerAllRoutes()
	s.Handler = s.router

	return s, nil
}

// Start launches the server.
func (s *Server) Start() error {
	log.Printf("Server listening at http://localhost%s\n", s.Addr)

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// registerAllRoutes registers each entity's routes on the server.
func (s *Server) registerAllRoutes() {
	s.router.HandleFunc("/", s.handleIndex)

	s.registerUserRoutes()
	s.registerAvatarRoutes()
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text-plain")
	w.WriteHeader(200)
	w.Write([]byte("Hello World!"))
}
