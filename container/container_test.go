package container

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/more-than-code/deploybot/util"
)

func TestImageBuild(t *testing.T) {

	localHelper := NewContainerHelper("unix:///var/run/docker.sock")

	path, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	repoName := "mo-service-graph"

	buf, err := util.TarFiles(fmt.Sprintf("%s/%s/", path, repoName))

	if err != nil {
		log.Fatalln(err)
	}

	tag := "binartist/" + repoName
	localHelper.BuildImage(buf, &types.ImageBuildOptions{Tags: []string{tag}})
}
