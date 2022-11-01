package task

import (
	"testing"

	"github.com/more-than-code/deploybot/model"
)

func TestPostInstall(t *testing.T) {
	shell := "go version"

	d := NewDeployer()
	err := d.Start(model.DeployConfig{PostInstall: shell})

	t.Log(err)
}
