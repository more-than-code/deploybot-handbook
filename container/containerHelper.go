package container

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type ContainerHelper struct {
	cli *client.Client
}

func NewContainerHelper(host string) *ContainerHelper {
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &ContainerHelper{cli: cli}
}

func (c *ContainerHelper) BuildImage(buildContext io.Reader, buidOptions *types.ImageBuildOptions) error {
	buildResponse, err := c.cli.ImageBuild(context.Background(), buildContext, *buidOptions)

	if err != nil {
		return err
	}

	res, err := io.ReadAll(buildResponse.Body)
	if err != nil {
		return err
	}

	log.Println(string(res))

	buildResponse.Body.Close()

	return nil
}

func (c *ContainerHelper) PushImage(imageTag string, pushOptions *types.ImagePushOptions) error {
	res, err := c.cli.ImagePush(context.Background(), imageTag, *pushOptions)

	if err != nil {
		return err
	}

	res.Close()

	return nil
}

func (c *ContainerHelper) StartContainer(imageName, containerName string) error {
	ctx := context.Background()

	reader, err := c.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, containerName)
	if err != nil {
		return err
	}

	if err := c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := c.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	var err2 error
	select {
	case err := <-errCh:
		if err != nil {
			log.Println(err)
			err2 = err
		}
	case <-statusCh:
	}

	out, err := c.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return err2
}
