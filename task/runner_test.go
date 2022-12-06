package task

import (
	"testing"

	"github.com/docker/docker/api/types/mount"
	"github.com/more-than-code/deploybot/model"
)

func TestBuildImage(t *testing.T) {
	r := NewRunner()
	err := r.DoTask(model.Task{Type: model.TaskBuild, Config: model.BuildConfig{
		ImageName: "binartist/mo-serive-graph",
		ImageTag:  "latest",
		RepoUrl:   "https://github.com/joe-and-his-friends/mo-service-graph.git",
		RepoName:  "mo-service-graph",
	}}, nil)

	if err != nil {
		t.Error(err)
	}
}

func TestRunContainer(t *testing.T) {
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
		"TOKEN_SECRET_KEY=123456",
	}

	r := NewRunner()
	err := r.DoTask(model.Task{Type: model.TaskDeploy, Config: model.DeployConfig{ExposedPort: "9000", HostPort: "9000", Env: env, ImageName: "binartist/mo-service-graph", ImageTag: "latest", Mounts: []mount.Mount{{Type: "bind", Source: "/var/opt/fcm", Target: "/var/opt/fcm"}}, ServiceName: "graph", AutoRemove: false}}, nil)

	if err != nil {
		t.Error(err)
	}

}
