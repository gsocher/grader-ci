package main

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/dpolansky/ci/server/amqp"
	"github.com/dpolansky/ci/worker"
)

// Sends sample build jobs to RabbitMQ.
func main() {
	client, err := amqp.NewClient("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to dial amqp")

	task := &worker.BuildTask{
		Language: "golang",
		CloneURL: "github.com/dpolansky/go-poet",
		ID:       uuid.New().String(),
	}
	b, err := json.Marshal(task)
	failOnError(err, "Failed to marshal task")

	err = client.SendToQueue("jobs", b)
	failOnError(err, "Failed to send message to queue")

	logrus.WithField("cloneURL", task.CloneURL).Infof("Sent task to queue")
}

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
	}
}
