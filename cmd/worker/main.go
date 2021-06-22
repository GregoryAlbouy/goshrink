package main

import (
	"flag"
	"log"

	"github.com/gregoryalbouy/goshrink/internal/database"
	"github.com/gregoryalbouy/goshrink/pkg/dotenv"
	"github.com/gregoryalbouy/goshrink/pkg/queue"
	"github.com/streadway/amqp"
)

const defaultEnvPath = "./.env"

var env = map[string]string{
	"MYSQL_USER":          "",
	"MYSQL_ROOT_PASSWORD": "",
	"MYSQL_HOST":          "",
	"MYSQL_PORT":          "",
	"MYSQL_DATABASE":      "",
	"QUEUE_URL":           "",
	"QUEUE_NAME":          "",
	"STORAGE_SERVER_KEY":  "",
	"STORAGE_SERVER_URL":  "",
}

func main() {
	width := flag.Int("w", 200, "resize width")
	flag.Parse()

	// Configure rezise width used by the worker
	resizeWidth = *width

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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
		Domain:   env["MYSQL_HOST"],
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
