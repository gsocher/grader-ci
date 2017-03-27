package service

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/dpolansky/ci/model"
)

type RepositoryReader interface {
	GetRepositoryByID(id int) (*model.Repository, error)
	GetRepositoriesByOwner(owner string) ([]*model.Repository, error)
	GetRepositories() ([]*model.Repository, error)
}

type RepositoryWriter interface {
	UpdateRepository(m *model.Repository) error
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

func (r *rep) UpdateRepository(m *model.Repository) error {
	ps := `INSERT OR REPLACE INTO repos (id, owner, name, avatar_url) values (?, ?, ?, ?)`
	stmt, err := r.db.Prepare(ps)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(m.ID, m.Owner, m.Name, m.AvatarURL)
	if err != nil {
		return err
	}

	return nil
}

func (r *rep) GetRepositories() ([]*model.Repository, error) {
	ps := `SELECT id, owner, name, avatar_url FROM repos`
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
		rows.Scan(&m.ID, &m.Owner, &m.Name, &m.AvatarURL)
		res = append(res, m)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res, nil
}

func (r *rep) GetRepositoriesByOwner(owner string) ([]*model.Repository, error) {
	ps := `SELECT id, owner, name, avatar_url FROM repos where owner = ?`
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
		rows.Scan(&m.ID, &m.Owner, &m.Name, &m.AvatarURL)
		res = append(res, m)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res, nil
}

func (r *rep) GetRepositoryByID(id int) (*model.Repository, error) {
	ps := `SELECT id, owner, name, avatar_url FROM repos WHERE id = ?`
	stmt, err := r.db.Prepare(ps)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	m := &model.Repository{}

	for rows.Next() {
		rows.Scan(&m.ID, &m.Owner, &m.Name, &m.AvatarURL)
		break
	}

	if m.Owner == "" {
		return nil, fmt.Errorf("No repository found with ID: %v", id)
	}

	return m, nil
}
