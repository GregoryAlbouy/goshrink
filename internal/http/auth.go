package http

import (
	"errors"
	"net/http"
	"strings"
)

const (
	workerApiKey = "hi_mom"
	bearer       = "Bearer "
)

func extractAPIKey(r *http.Request) (string, error) {
	v := r.Header.Get("Authorization")
	if !strings.HasPrefix(v, bearer) {
		return "", errors.New("invalid authorization headers")
	}
	return strings.TrimPrefix(v, bearer), nil
}

func authenticateAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := extractAPIKey(r)
		if err != nil {
			respondHTTPError(w, ErrUnauthorized.Wrap(err))
			return
		}
		if key != workerApiKey {
			respondHTTPError(w, ErrUnauthorized.Wrap(errors.New("invalid api key")))
			return
		}
		next.ServeHTTP(w, r)
	})
}
