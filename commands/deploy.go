package commands

import (
	"github.com/remyLemeunier/contactkey/utils"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the service in an environment",
}

type Deploy struct {
	Env     string
	Service string
}

func (d Deploy) execute() {
}

func (d Deploy) fill(config *utils.Config, service string, env string) {
	d.Env = env
	d.Service = service
}
