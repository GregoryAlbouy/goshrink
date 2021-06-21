package main

import (
	"log"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
	"github.com/GregoryAlbouy/shrinker/pkg/httplog"
)

const defaultEnvPath = "./.env"

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
	// GET /storage/<filename>
	http.Handle("/storage/", httplog.RequestLogger(handleFileServe("/storage", fs)))
	// POST /storage/avatar
	http.Handle("/storage/avatar", httplog.RequestLogger(requireAPIKey(handleImageUpload)))

	addr := ":" + env["STATIC_SERVER_PORT"]

	log.Printf("Static server listening at http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
