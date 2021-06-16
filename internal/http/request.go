package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// extractID retreives the route parameter "id" from the mux route variables.
// Also validates the ID is an integer.
func extractID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return id, nil
}

// decodeBody reads the given request body and writes the decoded data to dest.
// The body is expected to be encoded as JSON.
func decodeBody(body io.ReadCloser, dest interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return ErrUnprocessableEntity.Wrap(err)
	}
	return nil
}
