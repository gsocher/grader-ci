package main

import (
	"database/sql"

	"github.com/dpolansky/ci/model"
)

func main() {
	db, err := sql.Open("sqlite3", model.SQLiteFilepath)
	must(err)

	must(createBuildsTable(db))
}

func createBuildsTable(db *sql.DB) error {
	table := `
	CREATE TABLE IF NOT EXISTS builds(
		id INTEGER NOT NULL PRIMARY KEY,
		
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
