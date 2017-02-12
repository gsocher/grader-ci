package worker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/dpolansky/ci/model"
)

// DockerClient is an interface for docker client functionality
type DockerClient interface {
	StartContainer(image string, name string) (string, error)
	RunBuild(containerID, pathToBuildScript string, build *model.BuildStatus, wr io.Writer) error
	StopContainer(containerID string) error
}

// dClient wraps docker's client
type dClient struct {
	client *client.Client
}

func newDockerClient() (DockerClient, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &dClient{
		client: c,
	}, nil
}

// StartContainer creates and starts a container with the given name and image.
func (d *dClient) StartContainer(image string, name string) (string, error) {
	containerCfg := &container.Config{
		AttachStderr: true,
		AttachStdout: true,
		Image:        image,
		OpenStdin:    true,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
		User:         "ci",
	}

	container, err := d.client.ContainerCreate(context.Background(), containerCfg, nil, nil, name)
	if err != nil {
		return "", err
	}

	err = d.client.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return container.ID, nil
}

func (d *dClient) RunBuild(containerID, pathToBuildScript string, build *model.BuildStatus, wr io.Writer) error {
	err := d.copyToContainer(pathToBuildScript, containerID, "/home/ci")
	if err != nil {
		return err
	}

	execCfg := types.ExecConfig{
		Cmd:          []string{"bash", "-x", "/home/ci/build.sh", build.CloneURL},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	exec, err := d.client.ContainerExecCreate(context.Background(), containerID, execCfg)
	if err != nil {
		return err
	}

	resp, err := d.client.ContainerExecAttach(context.Background(), exec.ID, execCfg)
	if err != nil {
		return err
	}

	defer resp.Close()

	err = d.client.ContainerExecStart(context.Background(), exec.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, resp.Reader)
	if err != nil {
		return err
	}

	inspect, err := d.client.ContainerExecInspect(context.Background(), exec.ID)
	if err != nil {
		return err
	}

	if inspect.ExitCode != 0 {
		return fmt.Errorf("Build exited with non-zero exit code: %v", inspect.ExitCode)
	}

	return nil
}

func (d *dClient) copyToContainer(srcPath, containerID, dstPath string) error {
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

	return d.client.CopyToContainer(context.Background(), containerID, dstDir, preparedArchive, types.CopyToContainerOptions{})
}

func (d *dClient) StopContainer(containerID string) error {
	return d.client.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
		Force: true,
	})
}
