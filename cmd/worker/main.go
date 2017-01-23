package main

import (
	"context"
	"encoding/json"

	"github.com/dpolansky/ci/worker"
	"github.com/sirupsen/logrus"

	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log := logrus.New()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"jobs", // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	w, err := worker.NewWorker()
	failOnError(err, "Failed to create worker")

	log.Infof("Waiting for build tasks")
	go func() {
		for d := range msgs {
			var task worker.BuildTask

			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.WithError(err).Errorf("Failed to unmarshal task")
			}

			task.Ctx = context.Background()

			log.WithFields(logrus.Fields{
				"id":       task.ID,
				"cloneURL": task.CloneURL,
				"lang":     task.Language,
			}).Infof("Received task")

			w.RunBuild(&task, os.Stdout)
		}
	}()

	select {}
}
