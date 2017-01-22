package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type build struct {
	cloneURL string
	lang     string
}

func main() {
	// endpoint := "unix:///var/run/docker.sock"
	client, err := client.NewEnvClient()
	must(err)

	createdBody, err := client.ContainerCreate(context.Background(),
		&container.Config{
			AttachStderr: true,
			AttachStdout: true,
			Image:        "build-golang",
			OpenStdin:    true,
			Cmd:          []string{"/bin/bash"},
			Tty:          true,
		},
		nil, nil, "test-container")
	must(err)
	fmt.Printf("create resp: %v\n", createdBody)

	err = client.ContainerStart(context.Background(), "test-container", types.ContainerStartOptions{})
	must(err)

	execConfig := types.ExecConfig{
		Cmd:          []string{"echo", "hello"},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	idResp, err := client.ContainerExecCreate(context.Background(), "test-container", execConfig)
	must(err)

	_, err = client.ContainerExecInspect(context.Background(), idResp.ID)
	must(err)

	resp, err := client.ContainerExecAttach(context.Background(), idResp.ID, execConfig)

	err = client.ContainerExecStart(context.Background(), idResp.ID, types.ExecStartCheck{Detach: false, Tty: true})
	must(err)

	w, err := io.Copy(os.Stdout, resp.Reader)
	fmt.Printf("%v %v\n", w, err)

	dur := time.Duration(0)
	err = client.ContainerStop(context.Background(), "test-container", &dur)
	must(err)
	err = client.ContainerRemove(context.Background(), "test-container", types.ContainerRemoveOptions{})
	must(err)
}

func runBuild(b *build) {
	log.Printf("Running build url=%v lang=%v", b.cloneURL, b.lang)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
