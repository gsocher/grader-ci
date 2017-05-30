package main

import (
	"database/sql"
	"log"

	"github.com/dpolansky/grader-ci/backend/dbutil"
	"github.com/dpolansky/grader-ci/model"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Printf("Opening SQLite connection\n")
	conn, err := sql.Open("sqlite3", model.SQLiteFilepath)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	log.Printf("Creating tables\n")
	err = dbutil.CreateSQLiteTables(conn)
	if err != nil {
		panic(err)
	}
}
