package worker

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
)

type BuildTask struct {
	Language string
	CloneURL string
	ID       string
	Ctx      context.Context
}

type Worker struct {
	dockerClient *client.Client
}

// NewWorker constructs a new worker and initializes a docker client.
func NewWorker() (*Worker, error) {
	client, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Docker client initialization error: %v", err)
	}

	return &Worker{
		dockerClient: client,
	}, nil
}

// runBuild runs a given BuildTask and streams its output to a writer. Returns the exit code
// from the build if it runs successfully.
func (w *Worker) RunBuild(b *BuildTask, wr io.Writer) (int, error) {
	log := logrus.WithFields(logrus.Fields{
		"id": b.ID,
	})

	log.Info("Starting build")

	image, err := getImageForLanguage(b.Language)
	if err != nil {
		return 0, err
	}

	containerCfg := &container.Config{
		AttachStderr: true,
		AttachStdout: true,
		Image:        image,
		OpenStdin:    true,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
		User:         "ci",
	}

	log.Info("Creating container")
	container, err := w.dockerClient.ContainerCreate(b.Ctx, containerCfg, nil, nil, b.ID)
	if err != nil {
		return 0, err
	}

	defer w.cleanupBuild(b)

	log.Info("Starting container")
	err = w.dockerClient.ContainerStart(b.Ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return 0, err
	}

	path, err := getBuildScriptPathForLanguage(b.Language)
	if err != nil {
		return 0, err
	}

	log.WithField("path", path).Info("Copying build script to container")
	err = w.copyToContainer(path, b.ID, "/home/ci")
	if err != nil {
		return 0, err
	}

	execCfg := types.ExecConfig{
		Cmd:          []string{"bash", "-x", "/home/ci/build.sh", b.CloneURL},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	log.Info("Creating exec")
	exec, err := w.dockerClient.ContainerExecCreate(b.Ctx, container.ID, execCfg)
	if err != nil {
		return 0, err
	}

	log.WithField("execID", exec.ID).Info("Attaching to exec")
	resp, err := w.dockerClient.ContainerExecAttach(b.Ctx, exec.ID, execCfg)
	if err != nil {
		return 0, err
	}

	defer resp.Close()

	log.WithField("execID", exec.ID).Info("Starting exec")
	err = w.dockerClient.ContainerExecStart(b.Ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return 0, err
	}

	_, err = io.Copy(wr, resp.Reader)
	if err != nil {
		return 0, err
	}

	inspect, err := w.dockerClient.ContainerExecInspect(b.Ctx, exec.ID)
	if err != nil {
		return 0, err
	}

	timeout := time.Duration(0)
	log.Info("Stopping container")
	err = w.dockerClient.ContainerStop(b.Ctx, b.ID, &timeout)
	return inspect.ExitCode, err
}

func (w *Worker) cleanupBuild(b *BuildTask) {
	err := w.dockerClient.ContainerRemove(b.Ctx, b.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		logrus.WithField("id", b.ID).WithError(err).Error("Failed to cleanup build")
	}
}

func (w *Worker) copyToContainer(srcPath, containerID, dstPath string) error {
	srcInfo := archive.CopyInfo{
		Exists: true,
		IsDir:  false,
		Path:   srcPath,
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		panic(err)
	}
	defer srcArchive.Close()

	dstInfo := archive.CopyInfo{
		Exists: true,
		IsDir:  true,
		Path:   dstPath,
	}

	dstDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		panic(err)
	}

	defer preparedArchive.Close()

	return w.dockerClient.CopyToContainer(context.Background(), containerID, dstDir, preparedArchive, types.CopyToContainerOptions{})
}

func getImageForLanguage(language string) (string, error) {
	return "build-golang", nil
}

func getBuildScriptPathForLanguage(language string) (string, error) {
	return filepath.Abs("./build/golang/build.sh")
}
