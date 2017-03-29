package service

import "database/sql"

func execStatement(db *sql.DB, ps string, data ...interface{}) (int, error) {
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
