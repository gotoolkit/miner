package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
)

func NewClient() {
	cli, err := dockerclient.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}
