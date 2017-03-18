package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpolansky/ci/amqp"
	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/worker"
	"github.com/sirupsen/logrus"

	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log := logrus.New()

	client, err := amqp.NewAMQPDefaultClient()
	failOnError(err, "Failed to create AMQP client")

	w, err := worker.New()
	failOnError(err, "Failed to create worker")

	log.Infof("Waiting for builds")

	callback := func(body []byte) {
		var build model.BuildStatus
		err := json.Unmarshal(body, &build)
		if err != nil {
			log.WithError(err).Errorf("Failed to unmarshal build")
		}

		log.WithFields(logrus.Fields{
			"id":       build.ID,
			"cloneURL": build.CloneURL,
		}).Infof("Received build")

		// set status to running and send an update
		build.Status = model.StatusBuildRunning
		byt, _ := json.Marshal(&build)
		client.SendToQueue(model.AMQPStatusQueue, byt)

		buf := &bytes.Buffer{}
		io.Copy(os.Stdout, buf)

		if err := w.RunBuild(&build, buf); err != nil {
			log.WithFields(logrus.Fields{
				"id":       build.ID,
				"cloneURL": build.CloneURL,
			}).WithError(err).Errorf("Failed to run build")

			build.Status = model.StatusBuildFailed
			build.Log += fmt.Sprintf("Failed to run build: %v\n", err)
		} else {
			build.Status = model.StatusBuildPassed
		}

		build.Log += buf.String()

		// send passed/failed status and log
		byt, _ = json.Marshal(&build)
		client.SendToQueue(model.AMQPStatusQueue, byt)
	}

	die := make(chan struct{})
	client.ReadFromQueueWithCallback(model.AMQPBuildQueue, callback, die)

	select {}
}
