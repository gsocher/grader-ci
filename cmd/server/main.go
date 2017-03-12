package main

import (
	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server"
	"github.com/dpolansky/ci/server/repo"
	"github.com/dpolansky/ci/server/service"
	"github.com/sirupsen/logrus"
)

func main() {
	amqpClient, err := service.NewAMQPClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start AMQP client")
	}

	sqliteBuildRepo, err := repo.NewSQLiteBuildRepo(model.SQLiteFilepath)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start SQLite build repo")
	}

	builder := service.NewBuilder(amqpClient, sqliteBuildRepo)

	serv, err := server.New(builder)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start server")
	}
	serv.Serve()
}
