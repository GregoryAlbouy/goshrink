package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// extractID retreives the route parameter "id" from the mux route variables.
// Also validates the ID format.
func extractID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return id, nil
}
