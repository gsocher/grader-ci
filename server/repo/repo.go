package repo

import (
	"database/sql"
	"fmt"

	"github.com/dpolansky/ci/model"
)

type RepositoryRepo interface {
	CreateRepository(cloneURL string) (*model.Repository, error)
	GetRepositoryByCloneURL(cloneURL string) (*model.Repository, error)
}

type sqliteRepositoryRepo struct {
	db *sql.DB
}

func NewSQLiteRepositoryRepo(filePath string) (RepositoryRepo, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	return &sqliteRepositoryRepo{
		db: db,
	}, nil
}

func (s *sqliteRepositoryRepo) CreateRepository(cloneURL string) (*model.Repository, error) {
	ps := `INSERT INTO repos (clone_url) values (?)`
	stmt, err := s.db.Prepare(ps)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	_, err = stmt.Exec(cloneURL)
	if err != nil {
		return nil, err
	}

	m := &model.Repository{
		CloneURL: cloneURL,
	}

	return m, nil
}

func (s *sqliteRepositoryRepo) GetRepositoryByCloneURL(cloneURL string) (*model.Repository, error) {
	ps := `SELECT clone_url FROM repos WHERE clone_url = ?`
	stmt, err := s.db.Prepare(ps)
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
		rows.Scan(&m.CloneURL)
		break
	}

	if m.CloneURL == "" {
		return nil, fmt.Errorf("No repository found with cloneURL: %v", cloneURL)
	}

	return m, nil
}
