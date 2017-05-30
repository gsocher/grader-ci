package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/dpolansky/grader-ci/backend/dbutil"
	"github.com/dpolansky/grader-ci/model"
)

const testBuildServiceCleanupStatement = "delete from builds"

func TestInsertBuild(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testBuildServiceCleanupStatement, t)

	b, err := NewSQLiteBuildService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	status := &model.BuildStatus{
		ID:         1,
		Log:        "foo",
		LastUpdate: time.Time{},
		Status:     model.StatusBuildWaiting,
		Source: &model.RepositoryMetadata{
			ID:       0,
			Branch:   "master",
			CloneURL: "foo",
		},
		Test: &model.RepositoryMetadata{
			ID:       2,
			Branch:   "master",
			CloneURL: "bar",
		},
	}

	expected, err := b.UpdateBuild(status)
	if err != nil {
		t.Fatalf("failed to insert build: %v", err)
	}

	if expected.LastUpdate.IsZero() {
		t.Fatalf("expected build service to update timestamp")
	}

	// we should be able to fetch the build
	actual, err := b.GetBuildByID(expected.ID)
	if err != nil {
		t.Fatalf("failed to get build: %v", err)
	}

	// cant use deepequal on expected/actual because it checks pointer comparison
	// of source/test repositories, so have to check manually
	if !reflect.DeepEqual(expected.Source, actual.Source) ||
		!reflect.DeepEqual(expected.Test, actual.Test) ||
		expected.Tested != actual.Tested ||
		expected.ID != actual.ID ||
		expected.Log != actual.Log ||
		expected.Status != actual.Status ||
		!expected.LastUpdate.Equal(actual.LastUpdate) {
		t.Fatalf("expected %+v got %+v", expected, actual)
	}
}

func TestUpdateBuild(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testBuildServiceCleanupStatement, t)

	b, err := NewSQLiteBuildService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	status := &model.BuildStatus{
		ID:         1,
		Log:        "foo",
		LastUpdate: time.Time{},
		Status:     model.StatusBuildWaiting,
		Source: &model.RepositoryMetadata{
			ID:       0,
			Branch:   "master",
			CloneURL: "foo",
		},
		Test: &model.RepositoryMetadata{
			ID:       2,
			Branch:   "master",
			CloneURL: "bar",
		},
	}

	once, err := b.UpdateBuild(status)
	if err != nil {
		t.Fatalf("failed to insert build: %v", err)
	}

	before := once.LastUpdate

	twice, err := b.UpdateBuild(once)
	if err != nil {
		t.Fatalf("failed to update build: %v", err)
	}

	if twice.LastUpdate.Equal(before) {
		t.Fatalf("update should have updated timestamp")
	}
}

func TestGetBuildByIDEmpty(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testBuildServiceCleanupStatement, t)

	b, err := NewSQLiteBuildService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	if _, err := b.GetBuildByID(0); err == nil {
		t.Fatalf("expected err fetching build that doesnt exist")
	}
}

func TestGetBuildsByRepositoryID(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testBuildServiceCleanupStatement, t)

	b, err := NewSQLiteBuildService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	status := &model.BuildStatus{
		ID:         1,
		Log:        "foo",
		LastUpdate: time.Time{},
		Status:     model.StatusBuildWaiting,
		Source: &model.RepositoryMetadata{
			ID:       0,
			Branch:   "master",
			CloneURL: "foo",
		},
		Test: &model.RepositoryMetadata{
			ID:       2,
			Branch:   "master",
			CloneURL: "bar",
		},
	}

	_, err = b.UpdateBuild(status)
	if err != nil {
		t.Fatalf("failed to insert build: %v", err)
	}

	builds, err := b.GetBuildsByRepositoryID(status.Source.ID)
	if err != nil {
		t.Fatalf("failed to get builds by source id: %v", err)
	}

	l := len(builds)
	if l != 1 {
		t.Fatalf("wrong number of builds returned, got %v expected %v", l, 1)
	}
}
