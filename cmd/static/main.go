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
	"STATIC_FILE_PATH":   "",
	"STATIC_SERVER_PORT": "",
	"STATIC_SERVER_KEY":  "",
}

func main() {
	envPath := dotenv.GetPath(defaultEnvPath)
	if err := dotenv.Load(envPath, &env); err != nil {
		log.Fatal(err)
	}

	fs := http.Dir(env["STATIC_FILE_PATH"])

	// GET /static/<filename>
	http.Handle("/static/", handleFileServe("/static", fs))
	// POST /static/avatar
	http.HandleFunc("/static/avatar", authenticate(handleImageUpload))

	addr := ":" + env["STATIC_SERVER_PORT"]
	fmt.Printf("Server listening at http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

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
			http.Error(w, "404 page not found", 404)
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
	// Place the pointer back at the start of the file
	file.Seek(0, io.SeekStart)

	// Infer the extension
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
		http.Error(w, "internal error: failed to create file", 500)
		return
	}
	defer dst.Close()

	// Copy all bytes from the file to the destination on disk
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "internal error: failed to copy file", 500)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("Created\n"))
}

// bearer represents the string prefixing the authorization key contained in
// the authorization headers: "Bearer <key>".
const bearer = "Bearer "

func authenticate(next http.HandlerFunc) http.HandlerFunc {
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
