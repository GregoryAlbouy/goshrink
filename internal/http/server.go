package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gregoryalbouy/goshrink/internal"
	"github.com/gregoryalbouy/goshrink/pkg/httputil"
	"github.com/gregoryalbouy/goshrink/pkg/queue"
	"github.com/gregoryalbouy/goshrink/pkg/simplejwt"
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
func NewServer(addr string, repo Repository, qp queue.Producer, secretKey string) *Server {
	s := &Server{
		Server:     &http.Server{Addr: addr},
		Repository: repo,
		imageQueue: qp,
	}
	simplejwt.SetSecretKey([]byte(secretKey))

	return s
}

// Start launches the server.
func (s *Server) Start() error {
	s.initRouter()

	log.Printf("Server listening at http://localhost%s\n", s.Addr)

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) initRouter() {
	s.router = mux.NewRouter().StrictSlash(true)
	s.router.Use(httputil.RequestLogger)
	s.registerAllRoutes()
	s.Handler = s.router
}

// registerAllRoutes registers each entity's routes on the server.
func (s *Server) registerAllRoutes() {
	s.router.HandleFunc("/", s.handleIndex)

	s.registerAuthRoutes()
	s.registerUserRoutes()
	s.registerAvatarRoutes()
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text-plain")
	w.WriteHeader(200)
	w.Write([]byte("Hello World!"))
}
