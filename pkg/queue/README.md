# Queue

The package `queue` defines the queue's implementation. It exposes `Producer` and `Consumer` objects with methods to publish and listen messages on the queue.

`queue` is built to work with [RabbitMQ](https://www.rabbitmq.com/documentation.html).

To use a named queue, you must invoke and pass your queue name to `SetQueueName` before calling `Producer.Publish` or
`Consumer.Listen`.

## Get the connection

Whether you are using `queue` with a `Consumer` or a `Producer`, you first need to open a connection.

```go
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
if err != nil {
    log.Fatal("could not connect to message broker")
}
defer conn.Close()

// Optionally set a queue name, or use a randomized name if skipped.
queue.SetQueueName("task_queue")
```

## Publish a message to the queue with a `Producer`

```go
producer, err := queue.NewProducer(conn)
if err != nil {
    log.Fatal("could not connect to the queue")
}

if err := producer.Publish(message); err != nil {
    // handle error
}
```

## Listen to messages on the queue with a `Consumer`

```go
consumer, err := queue.NewConsumer(conn)
if err != nil {
    log.Fatal("could not connect to the queue")
}

if err := producer.Listen(); err != nil {
    // handle error
}
```

In order to be configurable and reusable, `Consumer.Listen` accepts a function as a parameter to define the job it is assigned.

```go
func (c *Consumer) Listen(func(d amqp.Delivery) error) error { /* */ }
```

Simply define the job and pass it to `Listen`:

```go
func configuredJob(d amqp.Delivery) error {
    log.Printf("Received a message: %s", d.Body)
    return nil
}

if err := producer.Listen(configuredJob); err != nil {
    // handle error
}
```
