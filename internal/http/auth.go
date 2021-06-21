package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/GregoryAlbouy/shrinker/pkg/crypto"
	"github.com/GregoryAlbouy/shrinker/pkg/simplejwt"
)

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   int    `json:"-"`
}

type ContextKey string

const (
	userKey      ContextKey = "user"
	bearerScheme string     = "Bearer "
)

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

	claims := simplejwt.NewClaims(fmt.Sprint(u.ID), time.Now().Add(12*time.Hour))
	token, err := simplejwt.NewSignedToken(claims)
	if err != nil {
		respondHTTPError(w, ErrInternal)
		return
	}
	respondJSON(w, 201, struct {
		Token string `json:"access_token"`
	}{token})
}

func (s *Server) requireAuth(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		if !strings.HasPrefix(a, bearerScheme) {
			respondHTTPError(w, ErrUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(a, bearerScheme)

		token, err := simplejwt.VerifiedToken(tokenString)
		if err != nil {
			respondHTTPError(w, ErrUnauthorized)
			return
		}

		idStr, err := simplejwt.ClaimsId(*token)
		if err != nil {
			respondHTTPError(w, ErrUnauthorized)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondHTTPError(w, ErrInternal)
			return
		}

		u, err := s.UserService.FindByID(id)
		if err != nil {
			respondHTTPError(w, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), userKey, &u)
		hf(w, r.WithContext(ctx))
	}
}

func userFromContext(ctx context.Context) *internal.User {
	if u := ctx.Value(userKey); u != nil {
		return u.(*internal.User)
	}
	return nil
}
