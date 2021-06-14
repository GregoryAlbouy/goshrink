package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// registerDialRoutes is a helper function for registering all avatar routes.
func (s *Server) registerAvatarRoutes(r *mux.Router) {
	r.HandleFunc("/users/{id:[0-9]+}/avatar", s.handleAvatarUpload).Methods("POST")
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 202, "Accepted")
}
