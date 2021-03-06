package service

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/dpolansky/grader-ci/pkg/backend/dbutil"
	"github.com/dpolansky/grader-ci/pkg/model"
	_ "github.com/mattn/go-sqlite3"
)

const (
	buildSelectAllQuery  = "select id, repo_id, repo_url, repo_branch, tested, test_id, test_url, test_branch, status, date, log from builds"
	buildInsertStatement = "insert into builds(repo_id, repo_url, repo_branch, tested, test_id, test_url, test_branch, status, date, log) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	buildUpdateStatement = "update builds SET date=?, status=?, log=? WHERE id=?"
)

type BuildService interface {
	GetBuildByID(id int) (*model.BuildStatus, error)
	GetBuildsByRepositoryID(id int) ([]*model.BuildStatus, error)
	GetBuilds() ([]*model.BuildStatus, error)
	UpdateBuild(m *model.BuildStatus) (*model.BuildStatus, error)
}

type builder struct {
	db *sql.DB
}

func NewSQLiteBuildService(db *sql.DB) (BuildService, error) {
	return &builder{
		db: db,
	}, nil
}

func (b *builder) UpdateBuild(m *model.BuildStatus) (*model.BuildStatus, error) {
	m.LastUpdate = time.Now()
	if m.Test == nil {
		m.Test = &model.RepositoryMetadata{}
	}

	if _, err := b.GetBuildByID(m.ID); err == nil {
		_, err = dbutil.ExecStatement(b.db, buildUpdateStatement, m.LastUpdate, m.Status, m.Log, m.ID)
		if err != nil {
			return nil, fmt.Errorf("Build update failed: %v", err)
		}
	} else {
		// could not find build, so insert it
		id, err := dbutil.ExecStatement(b.db, buildInsertStatement, m.Source.ID, m.Source.CloneURL, m.Source.Branch, m.Tested, m.Test.ID, m.Test.CloneURL, m.Test.Branch, m.Status, m.LastUpdate, m.Log)
		if err != nil {
			return nil, fmt.Errorf("Build insert failed: %v", err)
		}

		m.ID = id
	}

	return m, nil
}

func (b *builder) GetBuildByID(id int) (*model.BuildStatus, error) {
	res, err := queryBuilds(b.db, fmt.Sprintf("%s where id = ?", buildSelectAllQuery), id)
	if err != nil {
		return nil, err
	}

	l := len(res)
	if l == 0 {
		return nil, fmt.Errorf("No build found with id=%v", id)
	} else if l > 1 {
		return nil, fmt.Errorf("Got more than one build but expected one, len=%v res=%v", l, res)
	}

	return res[0], nil
}

func (b *builder) GetBuildsByRepositoryID(id int) ([]*model.BuildStatus, error) {
	return queryBuilds(b.db, fmt.Sprintf("%s where repo_id = ?", buildSelectAllQuery), id)
}

func (b *builder) GetBuilds() ([]*model.BuildStatus, error) {
	return queryBuilds(b.db, buildSelectAllQuery)
}

func queryBuilds(db *sql.DB, ps string, data ...interface{}) ([]*model.BuildStatus, error) {
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

	res := []*model.BuildStatus{}
	for rows.Next() {
		m := &model.BuildStatus{
			Source: &model.RepositoryMetadata{},
			Test:   &model.RepositoryMetadata{},
		}
		err = rows.Scan(&m.ID, &m.Source.ID, &m.Source.CloneURL, &m.Source.Branch, &m.Tested, &m.Test.ID, &m.Test.CloneURL, &m.Test.Branch, &m.Status, &m.LastUpdate, &m.Log)
		if err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	// higher build numbers appear first
	sort.Slice(res, func(i, j int) bool {
		return res[i].ID > res[j].ID
	})

	return res, nil
}
