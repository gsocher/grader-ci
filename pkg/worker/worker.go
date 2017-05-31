package worker

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/dpolansky/grader-ci/pkg/model"
)

// Build directory containing docker images/build scripts relative to GOPATH
const buildDir = "src/github.com/dpolansky/grader-ci/pkg/worker/build"

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
func (w *Worker) RunBuild(b *model.BuildStatus, wr io.Writer) (int, error) {
	id := fmt.Sprintf("%v", b.ID)

	sourceDir, err := cloneRepoIntoTempDir(id, b.Source.CloneURL, b.Source.Branch)
	if err != nil {
		return 0, fmt.Errorf("Failed to clone repo: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	// handle tested repos by cloning source & test repo, then copy test into source with rsync, removing test dir afterwards
	if b.Tested {
		tempDir, err := cloneRepoIntoTempDir(fmt.Sprintf("%v_test", id), b.Test.CloneURL, b.Test.Branch)
		if err != nil {
			return 0, fmt.Errorf("Failed to clone test repository for tested build: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// copy temp dir into source dir and override conflicts
		cmd := exec.Command("rsync", "-av", fmt.Sprintf("%v/", tempDir), sourceDir)
		if err = cmd.Run(); err != nil {
			return 0, fmt.Errorf("Failed to rsync test into source: %v", err)
		}
	}

	cfg, err := parseConfigInDir(sourceDir)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse config: %v", err)
	}

	image, err := getImageForLanguage(cfg.Language)
	if err != nil {
		return 0, fmt.Errorf("Failed to get image name: %v", err)
	}

	scriptPath, err := getBuildScriptPathForLanguage(cfg.Language)
	if err != nil {
		return 0, fmt.Errorf("Failed to get build script path: %v", err)
	}

	// if there are commands to run:
	// - read the setup script
	// - add the commands
	// - write it to a temp file to be copied to the container
	if cfg.Script != nil {
		f, err := ioutil.ReadFile(scriptPath)
		if err != nil {
			return 0, fmt.Errorf("Failed to read script from %v: %v", scriptPath, err)
		}

		buf := bytes.NewBuffer(f)
		for _, c := range cfg.Script {
			buf.WriteString(c)
		}

		temp, err := ioutil.TempFile("", "")
		if err != nil {
			return 0, fmt.Errorf("Failed to create temp file: %v", err)
		}

		defer func() {
			os.Remove(temp.Name())
			temp.Close()
		}()

		_, err = temp.Write(buf.Bytes())

		scriptPath = temp.Name()
	}

	containerName, err := w.dockerClient.StartContainer(image)

	// startContainer creates and starts the container, so if the start fails then
	// we defer a cleanup to remove the created container
	defer func() {
		err := w.dockerClient.StopContainer(containerName)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to stop container %v", containerName)
		}
	}()

	if err != nil {
		return 0, fmt.Errorf("Failed to start container for image %v: %v", image, err)
	}

	err = w.dockerClient.CopyToContainer(containerName, scriptPath, "/root/build.sh", false, false)
	if err != nil {
		panic(err)
	}

	w.dockerClient.CopyToContainer(containerName, sourceDir, "/root/", true, true)

	exit, err := w.dockerClient.RunBuild(containerName, b, filepath.Base(sourceDir), wr)
	if err != nil {
		return 0, fmt.Errorf("Failed to run build on container %v: %v", b.ID, err)
	}

	return exit, nil
}

// cloneRepoIntoTempDir clones the target repo into a temp dir and returns the path of the dir.
func cloneRepoIntoTempDir(id, cloneURL, branch string) (string, error) {
	path, err := ioutil.TempDir("", id)
	if err != nil {
		return "", fmt.Errorf("Failed to create temp dir: %v", err)
	}

	cmd := exec.Command("git", "clone", "-b", branch, cloneURL, path)

	b := &bytes.Buffer{}
	cmd.Stderr = b

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
