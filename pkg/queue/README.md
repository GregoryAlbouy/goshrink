# Queue

The package `queue` defines the queue's implementation. It exposes `Producer` and `Consumer` objects with methods to publish and listen messages on the queue.

These objects are called by other packages.

```go
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
if err != nil {
    log.Fatal("could not connect to message broker")
}
defer conn.Close()

producer, err := queue.NewProducer(conn)
if err != nil {
    log.Fatal("could not connect to the queue")
}
producer.Publish(message)
```

```go
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
if err != nil {
    log.Fatal("could not connect to message broker")
}
defer conn.Close()

consumer, err := queue.NewConsumer(conn)
if err != nil {
    log.Fatal("could not connect to the queue")
}
consumer.Listen()
```

In order to be configurable and reusable, `(*Consumer).Listen()` accepts a parameter defining the job it is assigned:

```go
func (c *Consumer) Listen(func(d amqp.Delivery) error) error { /* */ }

func configuredJob(d amqp.Delivery) error {
    log.Printf("Received a message: %s", d.Body)
    return nil
}

consumer := NewConsumer()
consumer.Listen(configuredJob)
```
