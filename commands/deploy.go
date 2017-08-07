package commands

import (
	"github.com/remyLemeunier/contactkey/context"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the service in an environment",
}

type Deploy struct {
	Env     string
	Service string
	Context *context.Context
}

func (d Deploy) execute() {
}

func (d Deploy) fill(context *context.Context, service string, env string) {
	d.Env = env
	d.Service = service
	d.Context = context
}
