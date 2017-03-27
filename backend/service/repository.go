package service

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/dpolansky/ci/model"
)

type RepositoryReader interface {
	GetRepositoryByCloneURL(cloneURL string) (*model.Repository, error)
	GetRepositoriesByOwner(owner string) ([]*model.Repository, error)
	GetRepositories() ([]*model.Repository, error)
}

type RepositoryWriter interface {
	CreateRepository(m *model.Repository) error
}

type RepositoryReadWriter interface {
	RepositoryReader
	RepositoryWriter
}

type rep struct {
	db *sql.DB
}

func NewSQLiteRepositoryReadWriter(db *sql.DB) (RepositoryReadWriter, error) {
	return &rep{
		db: db,
	}, nil
}

func (r *rep) CreateRepository(m *model.Repository) error {
	ps := `INSERT INTO repos (clone_url, owner) values (?, ?)`
	stmt, err := r.db.Prepare(ps)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(m.CloneURL, m.Owner)
	if err != nil {
		return err
	}

	return nil
}

func (r *rep) GetRepositories() ([]*model.Repository, error) {
	ps := `SELECT clone_url, owner FROM repos`
	stmt, err := r.db.Prepare(ps)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := []*model.Repository{}

	for rows.Next() {
		m := &model.Repository{}
		rows.Scan(&m.CloneURL, &m.Owner)
		res = append(res, m)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CloneURL < res[j].CloneURL
	})

	return res, nil
}

func (r *rep) GetRepositoriesByOwner(owner string) ([]*model.Repository, error) {
	ps := `SELECT clone_url, owner FROM repos where owner = ?`
	stmt, err := r.db.Prepare(ps)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(owner)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := []*model.Repository{}

	for rows.Next() {
		m := &model.Repository{}
		rows.Scan(&m.CloneURL, &m.Owner)
		res = append(res, m)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CloneURL < res[j].CloneURL
	})

	return res, nil
}

func (r *rep) GetRepositoryByCloneURL(cloneURL string) (*model.Repository, error) {
	ps := `SELECT clone_url, owner FROM repos WHERE clone_url = ?`
	stmt, err := r.db.Prepare(ps)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(cloneURL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	m := &model.Repository{}

	for rows.Next() {
		rows.Scan(&m.CloneURL, m.Owner)
		break
	}

	if m.CloneURL == "" {
		return nil, fmt.Errorf("No repository found with cloneURL: %v", cloneURL)
	}

	return m, nil
}
