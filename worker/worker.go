package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

// runBuild runs a given BuildTask and streams its output to a writer.
func (w *Worker) RunBuild(b *BuildTask, wr io.Writer) error {
	image, err := getImageForLanguage(b.Language)
	if err != nil {
		return err
	}

	containerCfg := &container.Config{
		AttachStderr: true,
		AttachStdout: true,
		Image:        image,
		OpenStdin:    true,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
	}

	container, err := w.dockerClient.ContainerCreate(b.Ctx, containerCfg, nil, nil, b.ID)
	if err != nil {
		return fmt.Errorf("Failed to create container for build id=%v, err=%v", b.ID, err)
	}

	defer w.cleanupBuild(b)

	err = w.dockerClient.ContainerStart(b.Ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("Failed to start container for build id=%v, err=%v", b.ID, err)
	}

	execCfg := types.ExecConfig{
		Cmd:          []string{"echo", "hello"},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	exec, err := w.dockerClient.ContainerExecCreate(b.Ctx, container.ID, execCfg)
	if err != nil {
		return fmt.Errorf("Failed to create exec for build id=%v, err=%v", b.ID, err)
	}

	resp, err := w.dockerClient.ContainerExecAttach(b.Ctx, exec.ID, execCfg)
	if err != nil {
		return fmt.Errorf("Failed to attach to exec for build id=%v, err=%v", b.ID, err)
	}

	defer resp.Close()

	err = w.dockerClient.ContainerExecStart(b.Ctx, exec.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return fmt.Errorf("Failed to start exec for build id=%v, err=%v", b.ID, err)
	}

	_, err = io.Copy(wr, resp.Reader)
	if err != nil {
		return fmt.Errorf("Failed to read from exec connection for build id=%v, err=%v", b.ID, err)
	}

	timeout := time.Duration(0)
	err = w.dockerClient.ContainerStop(b.Ctx, b.ID, &timeout)
	if err != nil {
		return fmt.Errorf("Failed to stop container for build id=%v, err=%v", b.ID, err)
	}

	return nil
}

func (w *Worker) cleanupBuild(b *BuildTask) {
	err := w.dockerClient.ContainerRemove(b.Ctx, b.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		log.Printf("Failed to cleanup build id=%v, err=%v", b.ID, err)
	}
}

func getImageForLanguage(language string) (string, error) {
	return "build-golang", nil
}
