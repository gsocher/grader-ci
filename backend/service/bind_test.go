package service

import (
	"os"
	"reflect"
	"testing"

	"github.com/dpolansky/grader-ci/backend/dbutil"
	"github.com/dpolansky/grader-ci/model"
)

const testTestBindsCleanupStatement = "delete from test_binds"

func TestMain(m *testing.M) {
	if err := dbutil.SetupTables(); err != nil {
		panic(err)
	}

	exit := m.Run()

	if err := dbutil.TeardownTables(); err != nil {
		panic(err)
	}

	os.Exit(exit)
}

func TestUpdateTestBind(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testTestBindsCleanupStatement, t)

	b, err := NewSQLiteTestBindService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	expected := &model.TestBind{
		SourceID:   0,
		TestID:     0,
		TestBranch: "master",
	}

	err = b.UpdateTestBind(expected)
	if err != nil {
		t.Fatalf("Failed to upsert test bind: %v", err)
	}

	actual, err := b.GetTestBindBySourceRepositoryID(expected.SourceID)
	if err != nil {
		t.Fatalf("Failed to get test bind by source id: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Get bind by source ID failed, expected %+v got %+v", expected, actual)
	}
}

func TestGetTestBindDoesntExist(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testTestBindsCleanupStatement, t)

	b, err := NewSQLiteTestBindService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	if _, err = b.GetTestBindBySourceRepositoryID(0); err == nil {
		t.Fatalf("expected err fetching bind that doesnt exist")
	}
}

func TestGetTestBinds(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testTestBindsCleanupStatement, t)

	b, err := NewSQLiteTestBindService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	expected := []*model.TestBind{
		&model.TestBind{
			SourceID:   0,
			TestID:     0,
			TestBranch: "master",
		},
		&model.TestBind{
			SourceID:   1,
			TestID:     0,
			TestBranch: "dev",
		},
		&model.TestBind{
			SourceID:   2,
			TestID:     1,
			TestBranch: "foo",
		},
	}

	// add each expected binding
	for _, m := range expected {
		if err := b.UpdateTestBind(m); err != nil {
			t.Fatalf("failed to update bind: %v", err)
		}
	}

	actual, err := b.GetTestBinds()
	if err != nil {
		t.Fatalf("failed to get test binds: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestGetTestBindsEmpty(t *testing.T) {
	conn := dbutil.SetupConnection(t)
	defer dbutil.TeardownConnection(conn, testTestBindsCleanupStatement, t)

	b, err := NewSQLiteTestBindService(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	binds, err := b.GetTestBinds()
	if err != nil {
		t.Fatalf("failed to get test binds: %v", err)
	}

	l := len(binds)
	if l != 0 {
		t.Fatalf("expected empty slice of test binds, got len %v", len(binds))
	}
}
