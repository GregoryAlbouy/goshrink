package http

import (
	"fmt"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/pkg/mimetype"
)

// registerAvatarRoutes is a helper function for registering all avatar routes.
func (s *Server) registerAvatarRoutes() {
	s.router.HandleFunc("/users/{id:[0-9]+}/avatar", s.handleAvatarUpload).Methods("POST")
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	file, headers, err := r.FormFile("upload")
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}
	defer file.Close()

	if !mimetype.IsImage(file) {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	id, err := extractID(r)
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	msg := []byte(fmt.Sprintf("user %d uploading %s", id, headers.Filename))
	if err = s.producer.Publish(msg); err != nil {
		respondHTTPError(w, ErrInternal.Wrap(err))
	}

	respondJSON(w, 202, "Accepted\n")
}
