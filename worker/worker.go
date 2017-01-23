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
	"github.com/docker/docker/pkg/archive"
)

type BuildTask struct {
	Language    string
	CloneURL    string
	ID          string
	BuildScript io.Reader
	Ctx         context.Context
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

	container, err := w.dockerClient.ContainerCreate(b.Ctx, containerCfg, nil, nil, b.ID)
	if err != nil {
		return 0, err
	}

	defer w.cleanupBuild(b)

	err = w.dockerClient.ContainerStart(b.Ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return 0, err
	}

	err = w.copyToContainer("build.sh", b.ID, "/home/ci")
	if err != nil {
		return 0, err
	}

	execCfg := types.ExecConfig{
		Cmd:          []string{"bash", "-x", "/home/ci/build.sh"},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	exec, err := w.dockerClient.ContainerExecCreate(b.Ctx, container.ID, execCfg)
	if err != nil {
		return 0, err
	}

	resp, err := w.dockerClient.ContainerExecAttach(b.Ctx, exec.ID, execCfg)
	if err != nil {
		return 0, err
	}

	defer resp.Close()

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
	err = w.dockerClient.ContainerStop(b.Ctx, b.ID, &timeout)
	return inspect.ExitCode, err
}

func (w *Worker) cleanupBuild(b *BuildTask) {
	err := w.dockerClient.ContainerRemove(b.Ctx, b.ID, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		log.Printf("Failed to cleanup build id=%v, err=%v", b.ID, err)
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
