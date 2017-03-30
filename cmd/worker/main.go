package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/dpolansky/ci/amqp"
	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/worker"

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

	die := make(chan struct{})
	client.ReadFromQueueWithCallback(model.AMQPBuildQueue, handleBuild(client, w), die)

	select {}
}

func handleBuild(client amqp.Messenger, w *worker.Worker) func([]byte) {
	return func(body []byte) {
		var build model.BuildStatus

		err := json.Unmarshal(body, &build)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to unmarshal build")
		}

		log := logrus.New().WithField("id", build.ID)
		b, _ := json.MarshalIndent(build, "", "\t")
		log.Infof("Running build:\n%v\n", string(b))

		// set status to running and send an update
		build.Status = model.StatusBuildRunning
		byt, _ := json.Marshal(&build)
		client.SendToQueue(model.AMQPStatusQueue, byt)

		buf := &bytes.Buffer{}
		io.Copy(os.Stdout, buf)

		if exit, err := w.RunBuild(&build, buf); err != nil {
			log.WithError(err).Errorf("Failed to run build")

			build.Status = model.StatusBuildError
			build.Log += fmt.Sprintf("\nFailed to run build: %v\n", err)
		} else if exit != 0 {
			build.Status = model.StatusBuildFailed
			build.Log += fmt.Sprintf("\nBuild exited with non-zero status code: %v\n", exit)
		} else {
			build.Status = model.StatusBuildPassed
		}

		build.Log += buf.String()

		// send passed/failed status and log
		byt, _ = json.Marshal(&build)
		client.SendToQueue(model.AMQPStatusQueue, byt)
	}
}
