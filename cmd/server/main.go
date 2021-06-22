package main

import (
	"flag"
	"log"

	"github.com/GregoryAlbouy/shrinker/internal/database"
	"github.com/GregoryAlbouy/shrinker/internal/http"
	"github.com/GregoryAlbouy/shrinker/mock"
	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
	"github.com/GregoryAlbouy/shrinker/pkg/queue"
	"github.com/GregoryAlbouy/shrinker/pkg/simplejwt"
	"github.com/streadway/amqp"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"API_SERVER_PORT":     "",
	"API_JWT_SECRET":      "",
	"QUEUE_URL":           "",
	"QUEUE_NAME":          "",
	"MYSQL_USER":          "",
	"MYSQL_ROOT_PASSWORD": "",
	"MYSQL_HOST":          "",
	"MYSQL_PORT":          "",
	"MYSQL_DATABASE":      "",
}

func main() {
	migrate := flag.Bool("m", false, "use mock users")
	flag.Parse()

	envPath := dotenv.GetPath(defaultEnvPath)

	if err := run(envPath, *migrate); err != nil {
		log.Fatal(err)
	}
}

func run(envPath string, migrate bool) error {
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
	qp, err := initQueue()
	if err != nil {
		log.Fatalf("rabbitmq error: %s", err)
	}
	defer qp.CloseConnection()

	queue.SetQueueName(env["QUEUE_NAME"])

	srv, err := initServer(db, qp)
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
		Domain:   env["MYSQL_HOST"],
		Port:     env["MYSQL_PORT"],
		Database: env["MYSQL_DATABASE"],
	}

	db.MustConnect(cfg)
	log.Printf("Server connected to database %s", env["MYSQL_DATABASE"])
	return db
}

// initQueue connects and initializes the queue.
func initQueue() (queue.Producer, error) {
	conn, err := amqp.Dial(env["QUEUE_URL"])
	if err != nil {
		return queue.Producer{}, err
	}

	queue.SetQueueName(env["QUEUE_NAME"])

	producer, err := queue.NewProducer(conn)
	if err != nil {
		return queue.Producer{}, err
	}

	log.Printf("Server connected to queue %s", env["QUEUE_NAME"])
	return producer, nil
}

func migrateMockUsers(db *database.DB) {
	userService := database.NewUserService(db)
	users, err := mock.GetUsersWithHashedPasswords()
	if err != nil {
		log.Fatal(err)
	}
	if err := userService.Migrate(users); err != nil {
		log.Fatal(err)
	}
}

func initServer(db *database.DB, qp queue.Producer) (*http.Server, error) {
	addr := ":" + env["API_SERVER_PORT"]
	repo := http.Repository{
		UserService: database.NewUserService(db),
	}

	simplejwt.SetSecretKey([]byte(env["API_JWT_SECRET"]))

	return http.NewServer(addr, repo, qp)
}
