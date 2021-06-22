package main

import (
	"fmt"
	"net/http"

	"github.com/gregoryalbouy/goshrink/pkg/httputil"
	"github.com/gregoryalbouy/goshrink/pkg/mimetype"
)

func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve user ID
	userID := r.FormValue("userId")
	if userID == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	// Retrieve the file
	file, headers, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure it is a valid image
	if !mimetype.IsImage(file) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Save user image file
	filepath, err := saveUniqueUserImage(userID, headers.Filename, file)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Return file URL
	fileURL := buildFileURL(filepath)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fileURL))
}

// buildFileURL returns a file URL string from the given filepath.
func buildFileURL(filepath string) string {
	return fmt.Sprintf("http://localhost:%s/%s", env["STORAGE_SERVER_PORT"], filepath)
}

func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := httputil.BearerToken(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if key != env["STORAGE_SERVER_KEY"] {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
