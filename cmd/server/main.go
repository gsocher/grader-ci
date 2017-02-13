package main

import (
	"github.com/dpolansky/ci/server"
	"github.com/dpolansky/ci/server/amqp"
	"github.com/sirupsen/logrus"
)

func main() {
	amqpClient, err := amqp.NewClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start AMQP client")
	}

	buildService := server.NewBuildService(amqpClient)

	serv, err := server.New(amqpClient, buildService)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start server")
	}
	serv.Serve()
}
