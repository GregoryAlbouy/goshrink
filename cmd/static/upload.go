package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/GregoryAlbouy/shrinker/pkg/httputil"
	"github.com/GregoryAlbouy/shrinker/pkg/mimetype"
)

// handleImageUpload handles requests to upload an image to the server.
// The uploaded image is save on the disk if the request is accepted.
func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

	filepath := buildFilepath(headers.Filename)

	// Create a destination on disk
	dst, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "internal error: failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy all bytes from the file to the destination on disk
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "internal error: failed to copy file", http.StatusInternalServerError)
		return
	}
	fileURL := buildFileURL(filepath)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fileURL))
}

// buildFilepath efficiently builds a filepath string from the given filename.
func buildFilepath(filename string) string {
	parts := []string{env["STATIC_FILE_PATH"], "/", filename}
	return buildString(parts...)
}

// buildFileURL efficiently builds a file URL string from the given filepath.
func buildFileURL(filepath string) string {
	parts := []string{"http://localhost", ":", env["STATIC_SERVER_PORT"], "/", filepath}
	return buildString(parts...)
}

func buildString(parts ...string) string {
	var sb strings.Builder
	for _, v := range parts {
		sb.WriteString(v)
	}
	return sb.String()
}

func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := httputil.BearerToken(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if key != env["STATIC_SERVER_KEY"] {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
