package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// extractID retreives the route parameter "id" from the mux route variables.
// Also validates the ID format.
func extractID(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 0)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return uint(id), nil
}
