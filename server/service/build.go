package repo

import (
	"database/sql"
	"fmt"

	"github.com/dpolansky/ci/model"
	_ "github.com/mattn/go-sqlite3"
)

type BuildReadWriter interface {
	BuildReader
	BuildWriter
}

type BuildReader interface {
	GetBuildByID(id int) (*model.BuildStatus, error)
	GetBuildsBySourceRepositoryURL(cloneURL string) ([]*model.BuildStatus, error)
	GetBuilds() ([]*model.BuildStatus, error)
}

type BuildWriter interface {
	UpdateBuild(m *model.BuildStatus) (int, error)
}

type buildService struct {
	db *sql.DB
}

func NewSQLiteBuildService(db *sql.DB) (BuildReadWriter, error) {
	return &buildService{
		db: db,
	}, nil
}

func (b *buildService) UpdateBuild(m *model.BuildStatus) (int, error) {
	// check if the build exists, if it does then we are updating
	if _, err := b.GetBuildByID(m.ID); err == nil {
		ps := `
		UPDATE builds SET date=?, status=?, log=? WHERE id=?
		`

		stmt, err := b.db.Prepare(ps)
		if err != nil {
			return 0, err
		}

		defer stmt.Close()
		stmt.Exec(m.LastUpdate, m.Status, m.Log, m.ID)
		return m.ID, nil
	}

	// build doesn't exist, create new one
	ps := `INSERT INTO builds(clone_url, branch, status, date, log) values (?, ?, ?, ?, ?)`
	stmt, err := b.db.Prepare(ps)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(m.CloneURL, m.Branch, m.Status, m.LastUpdate, "")
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (b *buildService) GetBuildByID(id int) (*model.BuildStatus, error) {
	ps := `SELECT id, clone_url, date, branch, log, status FROM builds WHERE id = ?`
	stmt, err := b.db.Prepare(ps)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	m := &model.BuildStatus{}

	for rows.Next() {
		err = rows.Scan(&m.ID, &m.CloneURL, &m.LastUpdate, &m.Branch, &m.Log, &m.Status)
		if err != nil {
			return nil, err
		}
		break
	}

	// no result found
	if m.CloneURL == "" {
		return nil, fmt.Errorf("No build found with ID: %v", id)
	}

	return m, nil
}

func (b *buildService) GetBuildsBySourceRepositoryURL(cloneURL string) ([]*model.BuildStatus, error) {
	ps := `SELECT id, clone_url, date, branch, log, status FROM builds WHERE clone_url = ?`
	stmt, err := b.db.Prepare(ps)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(cloneURL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := []*model.BuildStatus{}
	for rows.Next() {
		m := &model.BuildStatus{}
		err = rows.Scan(&m.ID, &m.CloneURL, &m.LastUpdate, &m.Branch, &m.Log, &m.Status)
		if err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	return res, nil
}

func (b *buildService) GetBuilds() ([]*model.BuildStatus, error) {
	ps := `SELECT id, clone_url, date, branch, log, status FROM builds`
	stmt, err := b.db.Prepare(ps)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := []*model.BuildStatus{}
	for rows.Next() {
		m := &model.BuildStatus{}
		err = rows.Scan(&m.ID, &m.CloneURL, &m.LastUpdate, &m.Branch, &m.Log, &m.Status)
		if err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	return res, nil
}
