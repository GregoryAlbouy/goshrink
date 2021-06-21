package authtoken

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

const (
	basic  = "Basic "
	bearer = "Bearer "
)

var ErrInvalidHeaders = errors.New("invalid authorization headers")

func BearerToken(r *http.Request) (string, error) {
	ah := r.Header.Get("Authorization")
	if !strings.HasPrefix(ah, bearer) {
		return "", ErrInvalidHeaders
	}
	return strings.TrimPrefix(ah, bearer), nil
}

func BasicToken(r *http.Request) (string, error) {
	ah := r.Header.Get("Authorization")
	if !strings.HasPrefix(ah, basic) {
		return "", ErrInvalidHeaders
	}
	str, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ah, basic))
	if err != nil {
		return "", ErrInvalidHeaders
	}
	creds := strings.Split(string(str), ":")
	return creds[0], nil
}
