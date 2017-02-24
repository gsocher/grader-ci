package worker

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"os"

	"github.com/dpolansky/ci/model"
	"github.com/sirupsen/logrus"
)

// Build directory containing docker images/build scripts relative to GOPATH
const buildDir = "src/github.com/dpolansky/ci/worker/build"

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

// RunBuild runs a given BuildTask and streams its output to a writer.
func (w *Worker) RunBuild(b *model.BuildStatus, wr io.Writer) error {
	log := logrus.WithFields(logrus.Fields{
		"id":       b.ID,
		"cloneURL": b.CloneURL,
		"branch":   b.Branch,
	})

	log.Infof("Initializing build")
	id := fmt.Sprintf("%v", b.ID)

	log.Infof("Cloning into repo")
	dir, err := cloneRepoIntoTempDir(id, b.CloneURL, b.Branch)
	if err != nil {
		return fmt.Errorf("Failed to clone repo: %v", err)
	}
	defer os.RemoveAll(dir)

	cfg, err := parseConfigInDir(dir)
	if err != nil {
		return fmt.Errorf("Failed to parse config: %v", err)
	}

	image, err := getImageForLanguage(cfg.Language)
	if err != nil {
		return fmt.Errorf("Failed to get image name: %v", err)
	}

	scriptPath, err := getBuildScriptPathForLanguage(cfg.Language)
	if err != nil {
		return fmt.Errorf("Failed to get build sript path: %v", err)
	}

	log.Info("Starting container")

	containerName, err := w.dockerClient.StartContainer(image)

	// startContainer creates and starts the container, so if the start fails then
	// we defer a cleanup to remove the created container
	defer func() {
		log.Info("Stopping container")
		err := w.dockerClient.StopContainer(containerName)
		if err != nil {
			log.WithError(err).Errorf("Failed to stop container")
		}
	}()

	if err != nil {
		return fmt.Errorf("Failed to start container for image %v: %v", image, err)
	}

	w.dockerClient.CopyToContainer(containerName, scriptPath, "/root", false)
	w.dockerClient.CopyToContainer(containerName, dir, "/root/", true)

	log.Info("Running build script")
	err = w.dockerClient.RunBuild(containerName, b, filepath.Base(dir), wr)
	if err != nil {
		return fmt.Errorf("Failed to run build on container %v: %v", b.ID, err)
	}

	return nil
}

// cloneRepoIntoTempDir clones the target repo into a temp dir and returns the path of the dir.
func cloneRepoIntoTempDir(id, cloneURL, branch string) (string, error) {
	path, err := ioutil.TempDir("", id)
	if err != nil {
		return "", fmt.Errorf("Failed to create temp dir: %v", err)
	}

	cmd := exec.Command("git", "clone", "-b", branch, cloneURL, path)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(path)
		return "", fmt.Errorf("Failed to exec clone command: %v", err)
	}

	return path, nil
}

// parseConfigInDir reads the config yaml file from a directory and returns
// a model.Config object representing the file.
func parseConfigInDir(path string) (*model.Config, error) {
	newPath := filepath.Join(path, model.ConfigFileName)
	b, err := ioutil.ReadFile(newPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %v", err)
	}

	return parse(b)
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
