package main

import (
	"log"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
	"github.com/GregoryAlbouy/shrinker/pkg/httputil"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"STORAGE_FILE_PATH":   "",
	"STORAGE_SERVER_PORT": "",
	"STORAGE_SERVER_KEY":  "",
}

func main() {
	envPath := dotenv.GetPath(defaultEnvPath)
	if err := dotenv.Load(envPath, &env); err != nil {
		log.Fatal(err)
	}

	router := initRouter()
	addr := ":" + env["STORAGE_SERVER_PORT"]
	log.Printf("Static server listening at http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, httputil.RequestLogger(router)); err != nil {
		log.Fatal(err)
	}
}

func initRouter() *http.ServeMux {
	router := http.NewServeMux()
	fs := http.Dir(env["STORAGE_FILE_PATH"])
	// GET /storage/<filename>
	router.Handle("/storage/", handleFileServe("/storage", fs))
	// POST /storage/avatar
	router.Handle("/storage/avatar", requireAPIKey(handleImageUpload))
	return router
}
