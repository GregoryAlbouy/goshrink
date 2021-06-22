package http

import (
	"errors"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/GregoryAlbouy/shrinker/pkg/crypto"
)

// registerUserRoutes is a helper function for registering all user routes.
func (s *Server) registerUserRoutes() {
	s.router.HandleFunc("/users/{username:[a-zA-Z0-9]+}", s.handleUserGet).Methods("GET")

	s.router.HandleFunc("/users", s.handleUserCreate).Methods("POST")
}

func (s *Server) handleUserGet(w http.ResponseWriter, r *http.Request) {
	username, err := extractRouteParam(r, "username")
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	u, err := s.UserService.FindByUsername(username)
	if err != nil {
		respondHTTPError(w, ErrNotFound)
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
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	// Do not store password as plain text.
	hashedPwd, err := crypto.HashPassword(u.Password)
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(errors.New("invalid password")))
		return
	}
	u.Password = hashedPwd

	if err := s.UserService.InsertOne(*u); err != nil {
		respondHTTPError(w, ErrInternal)
		return
	}
	respondJSON(w, 201, "Created\n")
}
