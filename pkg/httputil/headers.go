package httputil

import (
	"errors"
	"net/http"
	"strings"
)

const (
	// bearerScheme is the string prefixing the key or token
	// in the authorization headers: "Bearer <key>".
	bearerScheme = "Bearer "
)

func BearerToken(r *http.Request) (string, error) {
	ah := r.Header.Get("Authorization")
	if !strings.HasPrefix(ah, bearerScheme) {
		return "", errors.New("invalid authorization headers")
	}

	return strings.TrimPrefix(ah, bearerScheme), nil
}
