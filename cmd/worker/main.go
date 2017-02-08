package main

import (
	"encoding/json"

	"github.com/dpolansky/ci/worker"
	"github.com/sirupsen/logrus"

	"os"

	"github.com/dpolansky/ci/server/amqp"
	amqpAPI "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log := logrus.New()

	client, err := amqp.NewClient("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to create AMQP client")

	w, err := worker.New()
	failOnError(err, "Failed to create worker")

	log.Infof("Waiting for build tasks")

	callback := func(m amqpAPI.Delivery) {
		var task worker.BuildTask
		err := json.Unmarshal(m.Body, &task)
		if err != nil {
			log.WithError(err).Errorf("Failed to unmarshal task")
		}

		log.WithFields(logrus.Fields{
			"id":       task.ID,
			"cloneURL": task.CloneURL,
			"lang":     task.Language,
		}).Infof("Received task")

		w.RunBuild(&task, os.Stdout)
	}

	die := make(chan struct{})
	client.ReadFromQueueWithCallback("jobs", callback, die)

	select {}
}
