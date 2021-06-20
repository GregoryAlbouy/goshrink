package main

import (
	"fmt"
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
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load env
	envPath := dotenv.GetPath(defaultEnvPath)
	if err := dotenv.Load(envPath, &env); err != nil {
		return err
	}

	// init queue
	qc, err := initQueue()
	if err != nil {
		return err
	}
	defer qc.CloseConnection()

	// init databse
	db := initDatabase()
	defer db.Close()

	// handle queue messages
	h := queueHandler{
		userService: database.NewUserService(db),
	}
	return qc.Listen(h.handleMessage)
}

func initDatabase() *database.DB {
	db := &database.DB{}
	cfg := database.Config{
		User:     env["MYSQL_USER"],
		Password: env["MYSQL_ROOT_PASSWORD"],
		Domain:   env["MYSQL_DOMAIN"],
		Port:     env["MYSQL_PORT"],
		Database: env["MYSQL°DATABASE"],
	}

	db.MustConnect(cfg)
	fmt.Println("Worker successfully connected to database")
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

	fmt.Println("Worker successfully connected to queue")
	return consumer, nil
}
