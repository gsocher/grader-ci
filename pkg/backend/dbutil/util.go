package dbutil

import (
	"database/sql"
	"os"
	"testing"

	"fmt"

	"github.com/dpolansky/grader-ci/pkg/model"
)

func CreateSQLiteTables(db *sql.DB) error {
	tables := []func(*sql.DB) error{
		createReposTable,
		createBuildsTable,
		createTestBindsTable,
	}

	for _, f := range tables {
		err := f(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExecStatement(db *sql.DB, ps string, data ...interface{}) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(ps)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(data...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), tx.Commit()
}

func SetupTables() error {
	conn, err := sql.Open("sqlite3", model.SQLiteFilepath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite conn: %v", err)
	}

	if err = CreateSQLiteTables(conn); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close conn while setting up tables: %v", err)
	}

	return nil
}

func TeardownTables() error {
	if err := os.Remove(model.SQLiteFilepath); err != nil {
		return fmt.Errorf("failed to remove sqlite db: %v", err)
	}
	return nil
}

func SetupConnection(t *testing.T) *sql.DB {
	conn, err := sql.Open("sqlite3", model.SQLiteFilepath)
	if err != nil {
		t.Fatalf("failed to open sqlite conn: %v", err)
	}

	return conn
}

func TeardownConnection(conn *sql.DB, cleanupStatement string, t *testing.T) {
	if _, err := ExecStatement(conn, cleanupStatement); err != nil {
		t.Fatalf("failed to exec cleanup statement %v: %v", cleanupStatement, err)
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("failed to close sqlite conn: %v", err)
	}
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

func createTestBindsTable(db *sql.DB) error {
	table := `
	CREATE TABLE IF NOT EXISTS test_binds(
		source_id INTEGER NOT NULL PRIMARY KEY,
		test_id  INTEGER NOT NULL,
		test_branch TEXT NOT NULL,
		FOREIGN KEY (source_id) REFERENCES repos(id),
		FOREIGN KEY (test_id) REFERENCES repos(id)
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
		repo_url TEXT NOT NULL,
		repo_branch TEXT NOT NULL,
		tested boolean,
		test_id INTEGER,
		test_url TEXTL,
		test_branch TEXT,
		status text NOT NULL,
		date DATETIME NOT NULL,
		log text,
		FOREIGN KEY (repo_id) REFERENCES repos(id),
		FOREIGN KEY (test_id) REFERENCES repos(id)
	);
	`

	_, err := db.Exec(table)
	return err
}
