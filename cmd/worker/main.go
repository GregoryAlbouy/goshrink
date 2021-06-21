package main

import (
	"flag"
	"log"

	"github.com/GregoryAlbouy/shrinker/internal/database"
	"github.com/GregoryAlbouy/shrinker/pkg/dotenv"
	"github.com/GregoryAlbouy/shrinker/pkg/queue"
	"github.com/streadway/amqp"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"MYSQL_USER":          "",
	"MYSQL_ROOT_PASSWORD": "",
	"MYSQL_DOMAIN":        "",
	"MYSQL_PORT":          "",
	"MYSQL_DATABASE":      "",
	"QUEUE_URL":           "",
	"QUEUE_NAME":          "",
	"STATIC_SERVER_KEY":   "",
	"STATIC_SERVER_URL":   "",
}

func main() {
	w := flag.Int("w", 200, "resize width")
	flag.Parse()

	if err := run(*w); err != nil {
		log.Fatal(err)
	}
}

func run(w int) error {
	envPath := dotenv.GetPath(defaultEnvPath)
	if err := dotenv.Load(envPath, &env); err != nil {
		return err
	}

	qc, err := initQueue()
	if err != nil {
		return err
	}
	defer qc.CloseConnection()

	db := initDatabase()
	defer db.Close()

	// Configure rezise width used by the worker
	resizeWidth = w

	// Configure the message handler and start consuming the queue messages.
	h := messageHandler{
		userService: database.NewUserService(db),
	}
	return qc.Listen(h.handle)
}

func initDatabase() *database.DB {
	db := &database.DB{}
	cfg := database.Config{
		User:     env["MYSQL_USER"],
		Password: env["MYSQL_ROOT_PASSWORD"],
		Domain:   env["MYSQL_DOMAIN"],
		Port:     env["MYSQL_PORT"],
		Database: env["MYSQL_DATABASE"],
	}

	db.MustConnect(cfg)
	log.Printf("Worker connected to database %s", env["MYSQL_DATABASE"])
	return db
}

func initQueue() (queue.Consumer, error) {
	conn, err := amqp.Dial(env["QUEUE_URL"])
	if err != nil {
		return queue.Consumer{}, err
	}

	queue.SetQueueName(env["QUEUE_NAME"])

	consumer, err := queue.NewConsumer(conn)
	if err != nil {
		return queue.Consumer{}, err
	}

	log.Println("Worker connected to queue")
	return consumer, nil
}
