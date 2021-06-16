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
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	u, err := s.UserService.FindByID(id)
	if err != nil {
		respondHTTPError(w, ErrNotFound.Wrap(err))
		return
	}
	respondJSON(w, 200, u)
}

func (s *Server) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	in := internal.UserInput{}
	if err := decodeBody(r.Body, &in); err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	u := internal.NewUser(in)
	if err := u.Validate(); err != nil {
		respondHTTPError(w, ErrUnprocessableEntity.Wrap(err))
		return
	}

	if err := s.UserService.InsertOne(*u); err != nil {
		respondHTTPError(w, ErrInternal)
		return
	}
	respondJSON(w, 201, "Created\n")
}
