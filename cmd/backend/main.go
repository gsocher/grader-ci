package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dpolansky/grader-ci/pkg/amqp"
	"github.com/dpolansky/grader-ci/pkg/backend"
	"github.com/dpolansky/grader-ci/pkg/backend/service"
	"github.com/dpolansky/grader-ci/pkg/model"
	"github.com/sirupsen/logrus"
)

func main() {
	amqp, err := amqp.NewAMQPDefaultClient()
	must(err, "Failed to start AMQP client")

	db, err := sql.Open("sqlite3", model.SQLiteFilepath)
	must(err, "Failed to open SQLite database connection")
	defer db.Close()

	build, _ := service.NewSQLiteBuildService(db)
	rep, _ := service.NewSQLiteRepositoryService(db)
	bind, _ := service.NewSQLiteTestBindService(db)
	msg, _ := service.NewAMQPBuildMessageService(amqp, build)
	github := service.NewGithubWebhookService(&service.GithubWebhookServiceConfig{
		RepoService: rep,
		MsgService:  msg,
	})

	server, err := backend.New(&backend.ServerConfig{
		BuildMessageService:  msg,
		BuildService:         build,
		RepositoryService:    rep,
		TestBindService:      bind,
		GithubWebhookService: github,
	})
	must(err, "Failed to create server")

	server.Run()
}

func must(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatalf(msg)
	}
}
