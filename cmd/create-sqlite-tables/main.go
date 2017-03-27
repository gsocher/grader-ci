package main

import (
	"database/sql"
	"log"

	"github.com/dpolansky/ci/model"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Printf("Opening SQLite connection\n")
	db, err := sql.Open("sqlite3", model.SQLiteFilepath)
	must(err)

	log.Printf("Creating tables\n")
	must(createReposTable(db))
	must(createBuildsTable(db))
}

func createReposTable(db *sql.DB) error {
	table := `
	CREATE TABLE IF NOT EXISTS repos(
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		owner TEXT NOT NULL,
		avatar_url TEXT NOT NULL
	);
	`

	_, err := db.Exec(table)
	return err
}

func createBuildsTable(db *sql.DB) error {
	table := `
	CREATE TABLE IF NOT EXISTS builds(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		repo_id INTEGER NOT NULL,
		clone_url TEXT NOT NULL,
		date DATETIME NOT NULL,
		branch text NOT NULL,
		log text,
		status text,
		FOREIGN KEY (repo_id) REFERENCES repos(id)
	);
	`

	_, err := db.Exec(table)
	return err
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
