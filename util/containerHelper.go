package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot/model"
)

type ContainerHelperConfig struct {
	DhUsername string `envconfig:"DH_USERNAME"`
	DhPassword string `envconfig:"DH_PASSWORD"`
}

type ContainerHelper struct {
	cli *client.Client
	cfg ContainerHelperConfig
}

func NewContainerHelper(host string) *ContainerHelper {
	var cfg ContainerHelperConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &ContainerHelper{cli: cli, cfg: cfg}
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

func (h *ContainerHelper) PushImage(name string) error {
	authConfig := types.AuthConfig{
		Username: h.cfg.DhUsername,
		Password: h.cfg.DhPassword,
	}
	encodedJSON, _ := json.Marshal(authConfig)
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	res, err := h.cli.ImagePush(context.Background(), name, types.ImagePushOptions{RegistryAuth: authStr})

	if err != nil {
		return err
	}

	defer res.Close()
	io.Copy(os.Stdout, res)
	return nil
}

func (h *ContainerHelper) StartContainer(cfg *model.DeployConfig) {
	ctx := context.Background()

	h.cli.ContainerStop(ctx, cfg.ServiceName, nil)
	h.cli.ContainerRemove(ctx, cfg.ServiceName, types.ContainerRemoveOptions{})

	reader, err := h.cli.ImagePull(ctx, cfg.ImageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	cConfig := &container.Config{
		Image: cfg.ImageName,
		Env:   cfg.Env,
	}

	if cfg.ExposedPort != "" {
		cConfig.ExposedPorts = nat.PortSet{nat.Port(cfg.ExposedPort + "/tcp"): struct{}{}}
	}

	hConfig := &container.HostConfig{
		AutoRemove: cfg.AutoRemove,
	}

	if cfg.HostPort != "" {
		hConfig.PortBindings = nat.PortMap{nat.Port(cfg.HostPort + "/tcp"): []nat.PortBinding{{HostPort: cfg.HostPort, HostIP: "0.0.0.0"}}}
	}

	if cfg.MountSource != "" && cfg.MountTarget != "" {
		hConfig.Mounts = []mount.Mount{{Type: "bind", Source: cfg.MountSource, Target: cfg.MountTarget}}
	}

	resp, err := h.cli.ContainerCreate(ctx, cConfig, hConfig, &network.NetworkingConfig{}, nil, cfg.ServiceName)
	if err != nil {
		panic(err)
	}

	if err := h.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := h.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := h.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
