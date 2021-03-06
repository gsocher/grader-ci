package service

import (
	"database/sql"
	"fmt"

	"github.com/dpolansky/grader-ci/pkg/backend/dbutil"
	"github.com/dpolansky/grader-ci/pkg/model"
)

const (
	bindSelectAllQuery  = "select source_id, test_id, test_branch from test_binds"
	bindInsertStatement = "insert or replace into test_binds(source_id, test_id, test_branch) values (?, ?, ?)"
)

type TestBindService interface {
	GetTestBindBySourceRepositoryID(id int) (*model.TestBind, error)
	GetTestBinds() ([]*model.TestBind, error)
	UpdateTestBind(bind *model.TestBind) error
}

type binder struct {
	db *sql.DB
}

func NewSQLiteTestBindService(db *sql.DB) (TestBindService, error) {
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
	_, err := dbutil.ExecStatement(b.db, bindInsertStatement, bind.SourceID, bind.TestID, bind.TestBranch)
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
		err = rows.Scan(&m.SourceID, &m.TestID, &m.TestBranch)
		if err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	return res, nil
}
