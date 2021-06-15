package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/GregoryAlbouy/shrinker/internal"
	"github.com/GregoryAlbouy/shrinker/internal/database"
)

func migrateMockUsers(db *database.DB) {
	userService := database.NewUserService(db)

	users, err := readMockUsersFile(mockUsersPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := userService.Migrate(users); err != nil {
		log.Fatal(err)
	}
}

func readMockUsersFile(path string) ([]internal.User, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	users := []internal.User{}

	if err := json.Unmarshal(file, &users); err != nil {
		return nil, err
	}

	return users, nil
}
