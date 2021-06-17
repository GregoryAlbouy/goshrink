package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GregoryAlbouy/shrinker/internal/database"
	"github.com/GregoryAlbouy/shrinker/internal/http"
	"github.com/GregoryAlbouy/shrinker/mock"
	"github.com/joho/godotenv"
)

const (
	// defaultEnvPath is the path used to read environment variables,
	// if ENV_PATH is not set.
	defaultEnvPath = "./.env"
	mockUsersPath  = "./mock/users.json"
)

// env is a map of environment variables. It is set using loadEnv function.
var env = map[string]string{
	"API_SERVER_PORT":     "",
}

func main() {
	// Read migrate CLI flag
	migrate := flag.Bool("m", false, "use mock users")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()

	// Get environment file
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = defaultEnvPath
	}

	if err := run(envPath, *migrate, *verbose); err != nil {
		log.Fatal(err)
	}
}

func run(envPath string, migrate bool, verbose bool) error {
	if err := loadEnv(envPath); err != nil {
		return err
	}


// loadEnv reads values from the given filepath
// and stores them in `env` map.
// It returns an error if one is missing.
func loadEnv(filepath string) error {

	// Read env file
	envMap, err := godotenv.Read(filepath)
	if err != nil {
		return err
	}

	// Set env values and catch the missing ones
	missingEnv := []string{}
	for k := range env {
		if v, ok := envMap[k]; !ok {
			missingEnv = append(missingEnv, k)
		} else {
			env[k] = v
		}
	}

	// If one or more is missing, return an error
	// with a list of missing variables.
	if len(missingEnv) != 0 {
		missingEnvStr := strings.Join(missingEnv, ",")
		return fmt.Errorf("missing environment variables: %s", missingEnvStr)
	}

	return nil
}

func initServer(db *database.DB, verbose bool) *http.Server {
	addr := ":" + env["API_SERVER_PORT"]
	repo := http.Repository{
		UserService: database.NewUserService(db),
	}

	return http.NewServer(addr, repo, verbose)
}