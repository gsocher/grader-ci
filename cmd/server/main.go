package main

import (
	"github.com/docker/docker/api/server"
	"github.com/dpolansky/ci/backend/service"
	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/backend"
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

	sqliteRepositoryRepo, err := repo.NewSQLiteRepositoryRepo(model.SQLiteFilepath)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start SQlite repository repo")
	}

	builder := service.NewBuilder(amqpClient, sqliteBuildRepo)
	repositoryService := service.NewRepositoryService(sqliteRepositoryRepo)

	serv, err := server.New(builder, repositoryService)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start server")
	}
	serv.Serve()
}
