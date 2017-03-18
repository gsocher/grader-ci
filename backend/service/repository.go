package service

import (
	"database/sql"
	"fmt"

	"github.com/dpolansky/ci/model"
)

type RepositoryReader interface {
	GetRepositoryByCloneURL(cloneURL string) (*model.Repository, error)
	GetRepositoriesByOwner(owner string) ([]*model.Repository, error)
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
	return r.createRepositoryInDB(m)
}

func (r *rep) createRepositoryInDB(m *model.Repository) error {
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

func (r *rep) GetRepositoriesByOwner(owner string) ([]*model.Repository, error) {
	return r.getRepositoriesByOwnerInDB(owner)
}

func (r *rep) getRepositoriesByOwnerInDB(owner string) ([]*model.Repository, error) {
	ps := `SELECT clone_url, owner FROM repos order by clone_url asc`
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

	return res, nil
}

func (r *rep) GetRepositoryByCloneURL(cloneURL string) (*model.Repository, error) {
	return r.getRepositoryByCloneURLInDB(cloneURL)
}

func (r *rep) getRepositoryByCloneURLInDB(cloneURL string) (*model.Repository, error) {
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
