package util

import (
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/more-than-code/deploybot/model"
)

func TestBuildImage(t *testing.T) {

	localHelper := NewContainerHelper("unix:///var/run/docker.sock")

	path := "/var/opt/projects"

	repoName := "mo-service-graph"

	buf, err := TarFiles(fmt.Sprintf("%s/%s/", path, repoName))

	if err != nil {
		t.Error(err)
	}

	tag := "binartist/" + repoName
	err = localHelper.BuildImage(buf, &types.ImageBuildOptions{Tags: []string{tag}})

	if err != nil {
		t.Error(err)
	}
}

func TestPushImage(t *testing.T) {
	helper := NewContainerHelper("unix:///var/run/docker.sock")

	repoName := "mo-service-graph"
	tag := "binartist/" + repoName

	err := helper.PushImage(tag)

	if err != nil {
		t.Error(err)
	}

}

func TestRunImage(t *testing.T) {
	helper := NewContainerHelper("unix:///var/run/docker.sock")

	env := []string{
		"GRAPH_SERVER_PORT=:9000",
		"USER_SERVER_URI=localhost:8001",
		"SMS_SERVER_URI=localhost:8002",
		"AUTH_SERVER_URI=localhost:8003",
		"VOTING_SERVER_URI=localhost:8004",
		"MOMENT_SERVER_URI=localhost:8005",
		"WPAPI_SERVER_URI=localhost:8006",
		"TASK_SERVER_URI=localhost:8007",
		"REDIS_URI=test.mohiguide.com:6380",
		"OSS_ENDPOINT_DOWNLOAD=https://oss-mohiguide-com.sgp1.cdn.digitaloceanspaces.com",
		"OSS_ENDPOINT_UPLOAD=https://oss-mohiguide-com.sgp1.digitaloceanspaces.com",
		"OSS_REGION=sgp1",
		"OSS_SIGNING_SERVER=http://dev.mohiguide.com:8000",
		"WP_API_ENDPOINT_BASE=https://mh2.appsive.com/wp-json/wp/v2",
		"GIN_MODE=debug",
		"AT_TTL_MINUTE=0",
		"AT_TTL_HOUR=0",
		"AT_TTL_DAY=30",
		"RT_TTL_MINUTE=0",
		"RT_TTL_HOUR=0",
		"RT_TTL_DAY=30",
		"LATEST_RELEASED_VERSION=1.18.0",
		"LOWEST_SUPPORTED_VERSION=1.17.6",
		"DETAILS_URL=https://wordpress.uat.mohiguide.com/2022/06/02/hello-world/",
		"GOOGLE_APPLICATION_CREDENTIALS=/var/opt/fcm/app.json",
	}
	helper.StartContainer(&model.RunConfig{Env: env, ImageName: "binartist/mo-service-graph", ImageTag: "latest", Mounts: []mount.Mount{{Source: "/var/opt/fcm", Target: "/var/opt/fcm"}}, ServiceName: "graph", AutoRemove: true})
}
