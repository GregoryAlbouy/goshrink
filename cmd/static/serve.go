package main

import (
	"net/http"
	"strings"
)

// handleFileServe serves at the given path the files under the given directory.
func handleFileServe(path string, dir http.Dir) http.Handler {
	return http.StripPrefix(path, disableDirListing(http.FileServer(dir)))
}

// disableDirListing prevents http.FileServer from automatically generating
// navigable directory listings.
// It simply handles every path ending with a trailing slash as a 404.
func disableDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
