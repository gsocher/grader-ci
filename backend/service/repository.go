package service

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/dpolansky/grader-ci/backend/dbutil"
	"github.com/dpolansky/grader-ci/model"
)

const (
	repositorySelectAllQuery  = "select id, owner, name, avatar_url from repos"
	repositoryInsertStatement = "insert or replace into repos (id, owner, name, avatar_url) values (?, ?, ?, ?)"
)

type RepositoryService interface {
	GetRepositoryByID(id int) (*model.Repository, error)
	GetRepositoriesByOwner(owner string) ([]*model.Repository, error)
	GetRepositories() ([]*model.Repository, error)
	UpdateRepository(m *model.Repository) error
}

type rep struct {
	db *sql.DB
}

func NewSQLiteRepositoryService(db *sql.DB) (RepositoryService, error) {
	return &rep{
		db: db,
	}, nil
}

func (r *rep) UpdateRepository(m *model.Repository) error {
	_, err := dbutil.ExecStatement(r.db, repositoryInsertStatement, m.ID, m.Owner, m.Name, m.AvatarURL)
	return err
}

func (r *rep) GetRepositories() ([]*model.Repository, error) {
	return queryRepositories(r.db, repositorySelectAllQuery)
}

func (r *rep) GetRepositoriesByOwner(owner string) ([]*model.Repository, error) {
	return queryRepositories(r.db, fmt.Sprintf("%s where owner = ?", repositorySelectAllQuery), owner)
}

func (r *rep) GetRepositoryByID(id int) (*model.Repository, error) {
	res, err := queryRepositories(r.db, fmt.Sprintf("%s where id = ?", repositorySelectAllQuery), id)
	if err != nil {
		return nil, err
	}

	l := len(res)
	if l == 0 {
		return nil, fmt.Errorf("No repository found with id=%v", id)
	} else if l > 1 {
		return nil, fmt.Errorf("Got more than one repository but expected one, len=%v res=%v", l, res)
	}

	return res[0], nil
}

func queryRepositories(db *sql.DB, ps string, data ...interface{}) ([]*model.Repository, error) {
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
	res := []*model.Repository{}

	for rows.Next() {
		m := &model.Repository{}
		rows.Scan(&m.ID, &m.Owner, &m.Name, &m.AvatarURL)
		res = append(res, m)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res, nil
}
