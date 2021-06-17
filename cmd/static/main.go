package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
)

const (
	defaultEnvPath = "./.env"
)

// env is a map of environment variables. It is set using loadEnv function.
var env = map[string]string{
	"STATIC_SERVER_PORT": "",
	"STATIC_FILE_PATH":   "",
}

func main() {
	// Load env
	envPath := dotenv.GetPath(defaultEnvPath)
	dotenv.Load(envPath, &env)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./storage"))))
	http.HandleFunc("/static/avatar", handleFileUpload)

	addr := ":" + env["STATIC_SERVER_PORT"]
	fmt.Printf("Server listening at http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad request", 400)
		return
	}

	file, _, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	defer file.Close()

	f, err := os.Create("foo.png")
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("Created\n"))
}
