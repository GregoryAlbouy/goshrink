package http

import (
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/gorilla/mux"
)

// registerDialRoutes is a helper function for registering all user routes.
func (s *Server) registerUserRoutes(r *mux.Router) {
	r.HandleFunc("/users/{id:[0-9]+}", s.handleUserGet).Methods("GET")

	r.HandleFunc("/users", s.handleUserCreate).Methods("POST")

	r.HandleFunc("/users/{id:[0-9]+}/avatar", authenticateAPIKey(s.handleUserAvatarUpdate)).Methods("PATCH")
}

func (s *Server) handleUserGet(w http.ResponseWriter, r *http.Request) {
	id, _ := extractID(r)

	u := internal.User{
		ID:       id,
		Username: "string",
	}
	respondJSON(w, 200, u)
}

func (s *Server) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 201, "Created")
}

func (s *Server) handleUserAvatarUpdate(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 204, nil)
}
