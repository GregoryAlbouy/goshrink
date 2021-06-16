package http

import (
	"net/http"

	"github.com/GregoryAlbouy/shrinker/internal"
)

// registerAvatarRoutes is a helper function for registering all avatar routes.
func (s *Server) registerAvatarRoutes() {
	s.router.HandleFunc("/users/{id:[0-9]+}/avatar", s.handleAvatarUpload).Methods("POST")
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	in := internal.AvatarInput{}
	if err := decodeBody(r.Body, &in); err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	if err := in.Validate(); err != nil {
		respondHTTPError(w, ErrUnprocessableEntity.Wrap(err))
		return
	}

	if err := s.UserService.SetAvatarURL(id, in.URL); err != nil {
		respondHTTPError(w, ErrInternal)
		return
	}
	respondJSON(w, 202, "Accepted\n")
}
