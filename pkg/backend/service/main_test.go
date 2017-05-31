package service

import (
	"os"
	"testing"

	"github.com/dpolansky/grader-ci/pkg/backend/dbutil"
)

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
