package worker

import (
	"strings"
	"testing"

	"bytes"

	"github.com/dpolansky/grader-ci/pkg/model"
)

const testReposURL = "https://github.com/dpolansky/ci-test-repos"
const scriptTestBranchName = "golang-with-script"
const gradedSourceBranchName = "golang-graded-source"
const gradedTestBranchName = "golang-graded-test"

func TestGradedBuild(t *testing.T) {
	build := &model.BuildStatus{
		ID: 0,
		Source: &model.RepositoryMetadata{
			Branch:   gradedSourceBranchName,
			CloneURL: testReposURL,
		},
		Tested: true,
		Test: &model.RepositoryMetadata{
			Branch:   gradedTestBranchName,
			CloneURL: testReposURL,
		},
	}

	worker, err := New()
	if err != nil {
		t.Fatalf("failed to init worker: %v", err)
	}

	buf := &bytes.Buffer{}
	exit, err := worker.RunBuild(build, buf)
	output := buf.String()

	if err != nil {
		t.Fatalf("build error: %v", err)
	}

	if exit != 0 {
		t.Fatalf("build exited with non-zero: %v, output=%v", exit, output)
	}
}

func TestBuildWithScript(t *testing.T) {
	build := &model.BuildStatus{
		ID: 0,
		Source: &model.RepositoryMetadata{
			Branch:   scriptTestBranchName,
			CloneURL: testReposURL,
		},
	}

	worker, err := New()
	if err != nil {
		t.Fatalf("failed to init worker: %v", err)
	}

	buf := &bytes.Buffer{}
	exit, err := worker.RunBuild(build, buf)
	output := buf.String()

	if err != nil {
		t.Fatalf("build error: %v", err)
	}

	if exit != 0 {
		t.Fatalf("build exited with non-zero: %v, output=%v", exit, output)
	}

	if !strings.Contains(output, "done") {
		t.Fatalf("expected echo of 'done' but didn't see it, output=%v", output)
	}
}
