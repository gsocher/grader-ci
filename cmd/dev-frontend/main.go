package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dpolansky/grader-ci/backend/route"
	"github.com/dpolansky/grader-ci/backend/service"
	"github.com/dpolansky/grader-ci/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := sql.Open("sqlite3", model.SQLiteFilepath)
	must(err, "Failed to open SQLite database connection")
	defer db.Close()

	build, _ := service.NewSQLiteBuildReadWriter(db)
	rep, _ := service.NewSQLiteRepositoryReadWriter(db)
	bind, _ := service.NewSQLiteTestBindReadWriter(db)

	router := mux.NewRouter()
	route.RegisterRepositoryFrontendRoutes(router, rep)
	route.RegisterBuildFrontendRoutes(router, build, rep)
	route.RegisterBindFrontendRoutes(router, bind)
	route.RegisterAssetsRoute(router)

	serv := &http.Server{
		Addr:    "localhost:8081",
		Handler: router,
	}

	log.Printf("Listening on %v", serv.Addr)
	serv.ListenAndServe()
}

func must(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatalf(msg)
	}
}
