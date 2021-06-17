package http

import (
	"fmt"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/gorilla/mux"
)

type Server struct {
	*http.Server
	router *mux.Router
	Repository
}

// Repository exposes the available operations to access the data layer.
// Operations are grouped under services.
type Repository struct {
	UserService internal.UserService
}

// NewServer returns a new instance of Server given configuration parameters.
func NewServer(addr string, repo Repository) *Server {
	s := &Server{
		Server: &http.Server{Addr: addr},
		router: mux.NewRouter().StrictSlash(true),
		Repository: Repository{
			UserService: repo.UserService,
		},
	}

	s.registerAllRoutes()
	s.Handler = s.router

	return s
}

// Start launches the server.
func (s *Server) Start() error {
	fmt.Printf("Server listening at http://localhost%s\n", s.Addr)

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
