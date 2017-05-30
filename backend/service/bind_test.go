package service

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/dpolansky/grader-ci/backend/dbutil"
	"github.com/dpolansky/grader-ci/model"
)

var conn *sql.DB

func init() {
	var err error
	conn, err = sql.Open("sqlite3", model.SQLiteFilepath)
	if err != nil {
		panic(err)
	}

	if err = dbutil.CreateSQLiteTables(conn); err != nil {
		panic(err)
	}
}

func TestUpdate(t *testing.T) {
	binder, err := NewSQLiteTestBindReadWriter(conn)
	if err != nil {
		t.Fatalf("Failed to create bind service: %v", err)
	}

	expected := &model.TestBind{
		SourceID:   0,
		TestID:     0,
		TestBranch: "master",
	}

	err = binder.UpdateTestBind(expected)
	if err != nil {
		t.Fatalf("Failed to upsert test bind: %v", err)
	}

	actual, err := binder.GetTestBindBySourceRepositoryID(expected.SourceID)
	if err != nil {
		t.Fatalf("Failed to get test bind by source id: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Get bind by source ID failed, expected %+v got %+v", expected, actual)
	}
}
