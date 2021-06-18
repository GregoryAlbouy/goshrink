package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
	"github.com/GregoryAlbouy/shrinker/pkg/mimetype"
)

const (
	defaultEnvPath = "./.env"
)

var env = map[string]string{
	"STATIC_SERVER_PORT": "",
	"STATIC_FILE_PATH":   "",
}

func main() {
	envPath := dotenv.GetPath(defaultEnvPath)
	if err := dotenv.Load(envPath, &env); err != nil {
		log.Fatal(err)
	}

	// GET /static/<filename>
	http.Handle("/static/", handleFileServe("/static", env["STATIC_FILE_PATH"]))
	// POST /static/avatar
	http.HandleFunc("/static/avatar", handleImageUpload)

	addr := ":" + env["STATIC_SERVER_PORT"]
	fmt.Printf("Server listening at http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

// handleFileServe serves files under the given directory for the given path.
func handleFileServe(path string, dir string) http.Handler {
	return http.StripPrefix(path, disableDirListing(http.FileServer(http.Dir(dir))))
}

// disableDirListing prevents http.FileServer from automatically generating
// navigable directory listings.
// It simply handles every path ending with a trailing slash as a 404.
func disableDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Error(w, "not found", 404)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleImageUpload handles requests to upload an image to the server.
// The uploaded image is save on the disk if the request is accepted.
func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad request", 400)
		return
	}

	// Retrieve the file
	file, _, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	defer file.Close()

	// Ensure it is a valid image
	kind, err := mimetype.Detect(file)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	// Infer the extension.
	var ext string
	switch kind {
	case mimetype.PNG:
		ext = ".png"
	case mimetype.JPEG:
		ext = ".jpeg"
	default:
		http.Error(w, "bad request", 400)
		return
	}

	filepath := env["STATIC_FILE_PATH"] + "/foo" + ext

	// Create a destination on disk
	dst, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	defer dst.Close()

	// Copy all bytes from the file to the destination on disk
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("Created\n"))
}
