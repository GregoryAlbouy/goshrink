package main

import (
	"io"
	"net/http"
	"os"
	"strings"

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
	file, headers, err := r.FormFile("upload")
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
	// Place the pointer back at the start of the file
	file.Seek(0, io.SeekStart)

	// Create a destination on disk
	dst, err := os.OpenFile(getFilepath(headers.Filename), os.O_WRONLY|os.O_CREATE, 0666)
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

	w.WriteHeader(201)
	w.Write([]byte("Created\n"))
}

// getFilepath efficiently builds a filepath string from the given filename.
func getFilepath(filename string) string {
	var sb strings.Builder
	parts := []string{env["STATIC_FILE_PATH"], "/", filename}
	for _, v := range parts {
		sb.WriteString(v)
	}
	return sb.String()
}

// bearer represents the string prefixing the authorization key contained in
// the authorization headers: "Bearer <key>".
const bearer = "Bearer "

func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Header.Get("Authorization")
		if !strings.HasPrefix(v, bearer) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		key := strings.TrimPrefix(v, bearer)

		if key != env["STATIC_SERVER_KEY"] {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
