package main

import (
	"log"

	"github.com/GregoryAlbouy/shrinker/internal/http"
)

func main() {
	dbRepo := http.Repo{
		UserService:   nil,
		AvatarService: nil,
	}

	srv := http.NewServer(":9999", dbRepo)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
