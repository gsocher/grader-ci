package main

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/dpolansky/ci/model"
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

	log.Infof("Waiting for builds")

	callback := func(m amqpAPI.Delivery) {
		var build model.BuildStatus
		err := json.Unmarshal(m.Body, &build)
		if err != nil {
			log.WithError(err).Errorf("Failed to unmarshal build")
		}

		log.WithFields(logrus.Fields{
			"id":       build.ID,
			"cloneURL": build.CloneURL,
			"lang":     build.Language,
		}).Infof("Received build")

		// set status to running and send an update
		build.Status = model.StatusBuildRunning
		byt, _ := json.Marshal(&build)
		client.SendToQueue(model.AMQPStatusQueue, byt)

		buf := &bytes.Buffer{}
		io.Copy(os.Stdout, buf)

		if err := w.RunBuild(&build, buf); err != nil {
			build.Status = model.StatusBuildFailed
		} else {
			build.Status = model.StatusBuildPassed
		}

		build.Log = buf.String()

		// send passed/failed status and log
		byt, _ = json.Marshal(&build)
		client.SendToQueue(model.AMQPStatusQueue, byt)
	}

	die := make(chan struct{})
	client.ReadFromQueueWithCallback(model.AMQPBuildQueue, callback, die)

	select {}
}
