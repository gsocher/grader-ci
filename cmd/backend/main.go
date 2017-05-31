package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dpolansky/grader-ci/amqp"
	"github.com/dpolansky/grader-ci/backend"
	"github.com/dpolansky/grader-ci/backend/service"
	"github.com/dpolansky/grader-ci/model"
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

	serv, err := backend.New(build, msg, rep, bind)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to start server")
	}
	serv.Serve()
}

func must(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatalf(msg)
	}
}
