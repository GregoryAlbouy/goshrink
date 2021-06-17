package main

import (
	"flag"
	"log"

	"github.com/GregoryAlbouy/shrinker/internal/database"
	"github.com/GregoryAlbouy/shrinker/internal/http"
	"github.com/GregoryAlbouy/shrinker/mock"
	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
)

const (
	// defaultEnvPath is the path used to read environment variables,
	// if ENV_PATH is not set.
	defaultEnvPath = "./.env"
)

// env is a map of environment variables. It is set using loadEnv function.
var env = map[string]string{
	"API_SERVER_PORT":     "",
	"MYSQL_USER":          "",
	"MYSQL_ROOT_PASSWORD": "",
	"MYSQL_DOMAIN":        "",
	"MYSQL_PORT":          "",
	"MYSQL_DATABASE":      "",
}

func main() {
	// Read migrate CLI flag
	migrate := flag.Bool("m", false, "use mock users")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()

	// Get environment file
	envPath := dotenv.GetPath(defaultEnvPath)

	if err := run(envPath, *migrate, *verbose); err != nil {
		log.Fatal(err)
	}
}

func run(envPath string, migrate bool, verbose bool) error {
	if err := dotenv.Load(envPath, &env); err != nil {
		return err
	}

	db := mustInitDatabase()
	defer db.Close()

	// We use mock users for now for easier testing.
	// We might implement POST endpoints for that matter in the future.
	if migrate {
		migrateMockUsers(db)
	}

	srv := initServer(db, verbose)
	if err := srv.Start(); err != nil {
		return err
	}

	return nil
}

// mustInitDatabse connects and initializes the database.
// It panics if an error occurs.
func mustInitDatabase() *database.DB {
	db := &database.DB{}
	cfg := database.Config{
		User:     env["MYSQL_USER"],
		Password: env["MYSQL_ROOT_PASSWORD"],
		Domain:   env["MYSQL_DOMAIN"],
		Port:     env["MYSQL_PORT"],
		Database: env["MYSQL_DATABASE"],
	}

	db.MustInit(cfg)

	return db
}

func migrateMockUsers(db *database.DB) {
	userService := database.NewUserService(db)
	users := mock.Users
	if err := userService.Migrate(users); err != nil {
		log.Fatal(err)
	}
}

func initServer(db *database.DB, verbose bool) *http.Server {
	addr := ":" + env["API_SERVER_PORT"]
	repo := http.Repository{
		UserService: database.NewUserService(db),
	}

	return http.NewServer(addr, repo, verbose)
}
