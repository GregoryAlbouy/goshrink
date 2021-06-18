package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
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

	// GET /static/<filename>
	http.Handle("/static/", handleFileServe("/static", fs))
	// POST /static/avatar
	http.HandleFunc("/static/avatar", requireAPIKey(handleImageUpload))

	addr := ":" + env["STATIC_SERVER_PORT"]

	fmt.Printf("Server listening at http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
