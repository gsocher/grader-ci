package service

import (
	"database/sql"
	"fmt"

	"github.com/dpolansky/ci/model"
)

const (
	bindSelectAllQuery  = "select source_id, test_id from test_binds"
	bindInsertStatement = "insert or replace into test_binds(source_id, test_id) values (?, ?)"
)

type TestBindReader interface {
	GetTestBindBySourceRepositoryID(id int) (*model.TestBind, error)
	GetTestBinds() ([]*model.TestBind, error)
}

type TestBindWriter interface {
	UpdateTestBind(bind *model.TestBind) error
}

type TestBindReadWriter interface {
	TestBindReader
	TestBindWriter
}

type binder struct {
	db *sql.DB
}

func NewSQLiteTestBindReadWriter(db *sql.DB) (TestBindReadWriter, error) {
	return &binder{
		db: db,
	}, nil
}

func (b *binder) GetTestBindBySourceRepositoryID(id int) (*model.TestBind, error) {
	res, err := queryTestBinds(b.db, fmt.Sprintf("%s where source_id = ?", bindSelectAllQuery), id)
	if err != nil {
		return nil, err
	}

	l := len(res)
	if l == 0 {
		return nil, fmt.Errorf("No bind found with sourceID=%v", id)
	} else if l > 1 {
		return nil, fmt.Errorf("Got more than one bind but expected one, len=%v res=%v", l, res)
	}

	return res[0], nil
}

func (b *binder) GetTestBinds() ([]*model.TestBind, error) {
	return queryTestBinds(b.db, bindSelectAllQuery)
}

func (b *binder) UpdateTestBind(bind *model.TestBind) error {
	_, err := execStatement(b.db, bindInsertStatement, bind.SourceID, bind.TestID)
	return err
}

func queryTestBinds(db *sql.DB, ps string, data ...interface{}) ([]*model.TestBind, error) {
	stmt, err := db.Prepare(ps)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(data...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := []*model.TestBind{}
	for rows.Next() {
		m := &model.TestBind{}
		err = rows.Scan(&m.SourceID, &m.TestID)
		if err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	return res, nil
}
