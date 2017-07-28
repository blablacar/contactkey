package commands

import (
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff between what's currently deployed and what's going to be deployed",
}

type Diff struct {
	Env     string
	Service string
}

func (d Diff) execute() {

}

func (d Diff) fill(service string, env string) {
	d.Env = env
	d.Service = service
}
