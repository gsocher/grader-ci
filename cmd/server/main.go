package main

import (
	"encoding/json"

	"github.com/dpolansky/ci/worker"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

// Sends sample build jobs to RabbitMQ.
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"jobs", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	task := &worker.BuildTask{
		Language: "golang",
		CloneURL: "github.com/dpolansky/go-poet",
		ID:       uuid.New().String(),
	}

	bytes, err := json.Marshal(task)
	failOnError(err, "Failed to marshal build task")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(bytes),
		})

	logrus.WithFields(logrus.Fields{
		"id":       task.ID,
		"cloneURL": task.CloneURL,
		"lang":     task.Language,
	}).Infof("Sent task")
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
	}
}
