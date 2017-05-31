package service

import (
	"reflect"
	"testing"

	"github.com/dpolansky/grader-ci/pkg/backend/dbutil"
	"github.com/dpolansky/grader-ci/pkg/model"
)

const testRepositoryServiceCleanupStatement = "delete from repos"

func TestUpdateRepository(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testRepositoryServiceCleanupStatement, t)

	r, err := NewSQLiteRepositoryService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	expected := &model.Repository{
		AvatarURL: "avy",
		ID:        0,
		Name:      "foo",
		Owner:     "bar",
	}

	if err := r.UpdateRepository(expected); err != nil {
		t.Fatalf("failed to update repository: %v", err)
	}

	actual, err := r.GetRepositoryByID(expected.ID)
	if err != nil {
		t.Fatalf("failed to get repository by id: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestGetRepositoryByOwner(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testRepositoryServiceCleanupStatement, t)

	r, err := NewSQLiteRepositoryService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	expected := &model.Repository{
		AvatarURL: "avy",
		ID:        0,
		Name:      "foo",
		Owner:     "bar",
	}

	if err := r.UpdateRepository(expected); err != nil {
		t.Fatalf("failed to update repository: %v", err)
	}

	actual, err := r.GetRepositoriesByOwner(expected.Owner)
	if err != nil {
		t.Fatalf("failed to get repositories by owner: %v", err)
	}

	if !reflect.DeepEqual(expected, actual[0]) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestGetRepositoryByOwnerEmpty(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testRepositoryServiceCleanupStatement, t)

	r, err := NewSQLiteRepositoryService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	actual, err := r.GetRepositoriesByOwner("foo")
	if err != nil {
		t.Fatalf("failed to get repositories by owner: %v", err)
	}

	l := len(actual)
	if l != 0 {
		t.Fatalf("expected empty slice got %v", l)
	}
}

func TestGetRepositories(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testRepositoryServiceCleanupStatement, t)

	r, err := NewSQLiteRepositoryService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	expected := []*model.Repository{
		&model.Repository{
			AvatarURL: "avy",
			ID:        0,
			Name:      "foo",
			Owner:     "bar",
		},
		&model.Repository{
			AvatarURL: "avy",
			ID:        1,
			Name:      "a",
			Owner:     "b",
		},
		&model.Repository{
			AvatarURL: "avy",
			ID:        2,
			Name:      "c",
			Owner:     "d",
		},
	}

	for _, repo := range expected {
		if err := r.UpdateRepository(repo); err != nil {
			t.Fatalf("failed to update repository: %v", err)
		}
	}

	actual, err := r.GetRepositories()
	if err != nil {
		t.Fatalf("failed to get repositories: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestGetRepositoriesEmpty(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testRepositoryServiceCleanupStatement, t)

	r, err := NewSQLiteRepositoryService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	actual, err := r.GetRepositories()
	if err != nil {
		t.Fatalf("failed to get repositories by owner: %v", err)
	}

	l := len(actual)
	if l != 0 {
		t.Fatalf("expected empty slice got %v", l)
	}
}
