package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// extractRouteParam retreives the given route parameter from the
// mux path variables.
func extractRouteParam(r *http.Request, p string) (string, error) {
	v, ok := mux.Vars(r)[p]
	if !ok {
		return "", fmt.Errorf("invalid route parameter for \"%s\"", p)
	}
	return v, nil
}

// decodeBody reads the given request body and writes the decoded data to dest.
// The body is expected to be encoded as JSON.
func decodeBody(body io.ReadCloser, dest interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return err
	}
	return nil
}
