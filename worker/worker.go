package worker

import (
	"fmt"
	"io"
	"path/filepath"

	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Build directory containing docker images/build scripts relative to GOPATH
const buildDir = "src/github.com/dpolansky/ci/worker/build"

type BuildTask struct {
	Language string
	CloneURL string
	ID       string
}

type Worker struct {
	dockerClient DockerClient
}

// New constructs a new worker and initializes a docker client.
func New() (*Worker, error) {
	dockerClient, err := newDockerClient()
	if err != nil {
		return nil, err
	}

	return &Worker{
		dockerClient: dockerClient,
	}, nil
}

// runBuild runs a given BuildTask and streams its output to a writer.
func (w *Worker) RunBuild(b *BuildTask, wr io.Writer) error {
	b.ID = uuid.New().String()
	logrus.WithFields(logrus.Fields{
		"id":       b.ID,
		"lang":     b.Language,
		"cloneURL": b.CloneURL,
	}).Infof("Initializing build")

	image, err := getImageForLanguage(b.Language)
	if err != nil {
		return fmt.Errorf("Failed to get image name: %v", err)
	}

	pathToBuildScript, err := getBuildScriptPathForLanguage(b.Language)
	if err != nil {
		return fmt.Errorf("Failed to get build script for lang %v: %v", b.Language, err)
	}

	logrus.WithField("id", b.ID).Info("Starting container")
	_, err = w.dockerClient.StartContainer(image, b.ID)
	if err != nil {
		return fmt.Errorf("Failed to start container for image %v: %v", image, err)
	}

	defer func() {
		logrus.WithField("id", b.ID).Info("Stopping container")
		err := w.dockerClient.StopContainer(b.ID)
		if err != nil {
			logrus.WithField("id", b.ID).WithError(err).Errorf("Failed to stop container")
		}
	}()

	logrus.WithField("id", b.ID).WithField("scriptPath", pathToBuildScript).Info("Running build script")
	err = w.dockerClient.RunBuild(b.ID, pathToBuildScript, b, wr)
	if err != nil {
		return fmt.Errorf("Failed to run build on container %v: %v", b.ID, err)
	}

	return nil
}

// getImageForLanguage checks that a given language has a directory in build/ containing its
// build image and returns the name of the image.
func getImageForLanguage(language string) (string, error) {
	return fmt.Sprintf("build-%v", language), nil
}

func getBuildScriptPathForLanguage(language string) (string, error) {
	gopath := os.Getenv("GOPATH")
	path := filepath.Join(gopath, buildDir, language, "build.sh")
	return filepath.Abs(path)
}
