package http

import (
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
)

// registerUserRoutes is a helper function for registering all user routes.
func (s *Server) registerUserRoutes() {
	s.router.HandleFunc("/users/{id:[0-9]+}", s.handleUserGet).Methods("GET")

	s.router.HandleFunc("/users", s.handleUserCreate).Methods("POST")
}

func (s *Server) handleUserGet(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		respondHTTPError(w, ErrBadRequest)
		return
	}

	u := internal.User{
		ID:       id,
		Username: "string",
	}
	respondJSON(w, 200, u)
}

func (s *Server) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 201, "Created")
}
