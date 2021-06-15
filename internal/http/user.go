package http

import (
	"fmt"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
)

// registerUserRoutes is a helper function for registering all user routes.
func (s *Server) registerUserRoutes() {
	s.router.HandleFunc("/users/{id:[0-9]+}", s.handleUserGet).Methods("GET")

	s.router.HandleFunc("/users", s.handleUserCreate).Methods("POST")
}

func (s *Server) handleUserGet(w http.ResponseWriter, r *http.Request) {
	id, _ := extractID(r)

	u := internal.User{
		ID:       id,
		Username: "string",
	}
	serAvatarUrl(&u, r)
	respondJSON(w, 200, u)
}

func (s *Server) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 201, "Created")
}

func serAvatarUrl(u *internal.User, r *http.Request) {
	u.AvatarURL = fmt.Sprintf("%s/avatar", r.URL.String())
}
