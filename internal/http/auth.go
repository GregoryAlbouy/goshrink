package http

import (
	"net/http"

	"github.com/GregoryAlbouy/shrinker/pkg/crypto"
)

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) registerAuthRoutes() {
	s.router.HandleFunc("/login", s.handleLogin).Methods("POST")
}

// handleLogin handles the "GET /login" route. It simply renders an HTML login form.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	creds := &Creds{}
	if err := decodeBody(r.Body, creds); err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
	}

	u, err := s.UserService.FindCreds(creds.Username)
	if err != nil {
		respondHTTPError(w, ErrUnauthorized) // do not say whether or not this user exists
		return
	}

	if err := crypto.ComparePasswords(u.Password, creds.Password); err != nil {
		respondHTTPError(w, ErrUnauthorized)
		return
	}

	respondJSON(w, 201, "TOKEN")
}
