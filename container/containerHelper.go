package container

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/more-than-code/deploybot/model"
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

func (h *ContainerHelper) BuildImage(buildContext io.Reader, buidOptions *types.ImageBuildOptions) error {
	buildResponse, err := h.cli.ImageBuild(context.Background(), buildContext, *buidOptions)

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

func (h *ContainerHelper) PushImage(imageTag string, pushOptions *types.ImagePushOptions) error {
	res, err := h.cli.ImagePush(context.Background(), imageTag, *pushOptions)

	if err != nil {
		return err
	}

	res.Close()

	return nil
}

func (h *ContainerHelper) StartContainer(cfg *model.ContainerConfig) error {
	ctx := context.Background()

	reader, err := h.cli.ImagePull(ctx, cfg.ImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	resp, err := h.cli.ContainerCreate(ctx, &container.Config{
		Image: cfg.ImageName,
	}, &container.HostConfig{
		AutoRemove: cfg.AutoRemove,
		Mounts:     []mount.Mount{{Source: cfg.MountSource, Target: cfg.MountTarget, Type: mount.Type(cfg.MountType)}},
	}, nil, nil, cfg.ServiceName)
	if err != nil {
		return err
	}

	if err := h.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	// statusCh, errCh := h.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	// var err2 error
	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		log.Println(err)
	// 		err2 = err
	// 	}
	// case <-statusCh:
	// }

	// out, err := h.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	// if err != nil {
	// 	return err
	// }

	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}
