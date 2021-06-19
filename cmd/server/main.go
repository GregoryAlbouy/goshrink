package main

import (
	"flag"
	"log"

	"github.com/GregoryAlbouy/shrinker/internal/database"
	"github.com/GregoryAlbouy/shrinker/internal/http"
	"github.com/GregoryAlbouy/shrinker/mock"
	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
	"github.com/GregoryAlbouy/shrinker/pkg/queue"
	"github.com/streadway/amqp"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"API_SERVER_PORT":     "",
	"QUEUE_URL":           "",
	"QUEUE_NAME":          "",
	"MYSQL_USER":          "",
	"MYSQL_ROOT_PASSWORD": "",
	"MYSQL_DOMAIN":        "",
	"MYSQL_PORT":          "",
	"MYSQL_DATABASE":      "",
}

func main() {
	migrate := flag.Bool("m", false, "use mock users")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()

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

	// Connect to the queue as close to main as possible, as we are usign `defer`.
	q, err := amqp.Dial(env["QUEUE_URL"])
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()

	queue.SetQueueName(env["QUEUE_NAME"])

	srv, err := initServer(db, q, verbose)
	if err != nil {
		return err
	}

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

func initServer(db *database.DB, q *amqp.Connection, verbose bool) (*http.Server, error) {
	addr := ":" + env["API_SERVER_PORT"]
	repo := http.Repository{
		UserService: database.NewUserService(db),
	}

	return http.NewServer(addr, repo, q, verbose)
}
