package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dpolansky/ci/amqp"
	"github.com/dpolansky/ci/backend"
	"github.com/dpolansky/ci/backend/service"
	"github.com/dpolansky/ci/model"
	"github.com/sirupsen/logrus"
)

func main() {
	amqp, err := amqp.NewAMQPDefaultClient()
	must(err, "Failed to start AMQP client")

	db, err := sql.Open("sqlite3", model.SQLiteFilepath)
	must(err, "Failed to open SQLite database connection")
	defer db.Close()

	build, _ := service.NewSQLiteBuildReadWriter(db)
	rep, _ := service.NewSQLiteRepositoryReadWriter(db)
	bind, _ := service.NewSQLiteTestBindReadWriter(db)
	run, _ := service.NewAMQPBuildRunner(amqp, build)

	serv, err := backend.New(build, run, rep, bind)
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
