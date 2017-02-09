package main

import (
	"github.com/dpolansky/ci/server"
	"github.com/sirupsen/logrus"
)

func main() {
	serv, err := server.New()
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start server")
	}
	serv.Serve()
}
