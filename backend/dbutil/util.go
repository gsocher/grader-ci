package dbutil

import "database/sql"

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
	stmt, err := db.Prepare(ps)
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

	return int(id), nil
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
