package http

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/pkg/mimetype"
)

// registerAvatarRoutes is a helper function for registering all avatar routes.
func (s *Server) registerAvatarRoutes() {
	s.router.HandleFunc("/avatar", s.requireAuth(s.handleAvatarUpload)).Methods("POST")
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r.Context())

	file, _, err := r.FormFile("image")
	if err != nil {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}
	defer file.Close()

	if !mimetype.IsImage(file) {
		respondHTTPError(w, ErrBadRequest.Wrap(err))
		return
	}

	msg := new(bytes.Buffer)
	if _, err := msg.ReadFrom(file); err != nil {
		respondHTTPError(w, ErrInternal.Wrap(err))
		return
	}

	if err = s.imageQueue.Publish(msg.Bytes(), fmt.Sprint(user.ID)); err != nil {
		respondHTTPError(w, ErrInternal.Wrap(err))
		return
	}

	respondJSON(w, 202, "Accepted")
}
