package http

import (
	"net/http"

	"github.com/GregoryAlbouy/shrinker/pkg/mimetype"
)

// registerAvatarRoutes is a helper function for registering all avatar routes.
func (s *Server) registerAvatarRoutes() {
	s.router.HandleFunc("/users/{id:[0-9]+}/avatar", s.handleAvatarUpload).Methods("POST")
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("upload")
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}
	defer file.Close()

	if !mimetype.IsImage(file) {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	// There is no use for the user id for now. It will be passed to the queue.
	_, err = extractID(r)
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	respondJSON(w, 202, "Accepted\n")
}
