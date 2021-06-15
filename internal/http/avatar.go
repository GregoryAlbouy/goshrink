package http

import (
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
)

// registerAvatarRoutes is a helper function for registering all avatar routes.
func (s *Server) registerAvatarRoutes() {
	s.router.HandleFunc("/users/{id:[0-9]+}/avatar", s.handleAvatarGet).Methods("GET")

	s.router.HandleFunc("/users/{id:[0-9]+}/avatar", s.handleAvatarUpload).Methods("POST")
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 202, "Accepted")
}

func (s *Server) handleAvatarGet(w http.ResponseWriter, r *http.Request) {
	id, _ := extractID(r)

	a := internal.Avatar{
		UserID: id,
	}
	respondJSON(w, 200, a)
}
